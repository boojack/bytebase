package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/type/expr"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

var (
	testTenantNumber = 1
	prodTenantNumber = 1
	testInstanceName = "testInstanceTest"
	prodInstanceName = "testInstanceProd"
)

func TestDatabaseGroup(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	ctx := context.Background()
	ctl := &controller{}
	ctx, err := ctl.StartServerWithExternalPg(ctx)
	a.NoError(err)
	defer ctl.Close(ctx)

	// Create a project.
	projectID := generateRandomString("project", 10)
	project, err := ctl.projectServiceClient.CreateProject(ctx, &v1pb.CreateProjectRequest{
		Project: &v1pb.Project{
			Name:  fmt.Sprintf("projects/%s", projectID),
			Title: projectID,
		},
		ProjectId: projectID,
	})
	a.NoError(err)

	// Provision instances.
	instanceRootDir := t.TempDir()

	var testInstanceDirs []string
	var prodInstanceDirs []string
	for i := 0; i < testTenantNumber; i++ {
		instanceDir, err := ctl.provisionSQLiteInstance(instanceRootDir, fmt.Sprintf("%s-%d", testInstanceName, i))
		a.NoError(err)
		testInstanceDirs = append(testInstanceDirs, instanceDir)
	}
	for i := 0; i < prodTenantNumber; i++ {
		instanceDir, err := ctl.provisionSQLiteInstance(instanceRootDir, fmt.Sprintf("%s-%d", prodInstanceName, i))
		a.NoError(err)
		prodInstanceDirs = append(prodInstanceDirs, instanceDir)
	}

	// Add the provisioned instances.
	var testInstances []*v1pb.Instance
	var prodInstances []*v1pb.Instance
	for i, testInstanceDir := range testInstanceDirs {
		instance, err := ctl.instanceServiceClient.CreateInstance(ctx, &v1pb.CreateInstanceRequest{
			InstanceId: generateRandomString("instance", 10),
			Instance: &v1pb.Instance{
				Title:       fmt.Sprintf("%s-%d", testInstanceName, i),
				Engine:      v1pb.Engine_SQLITE,
				Environment: "environments/test",
				Activation:  true,
				DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: testInstanceDir, Id: "admin"}},
			},
		})
		a.NoError(err)
		testInstances = append(testInstances, instance)
	}
	for i, prodInstanceDir := range prodInstanceDirs {
		instance, err := ctl.instanceServiceClient.CreateInstance(ctx, &v1pb.CreateInstanceRequest{
			InstanceId: generateRandomString("instance", 10),
			Instance: &v1pb.Instance{
				Title:       fmt.Sprintf("%s-%d", prodInstanceName, i),
				Engine:      v1pb.Engine_SQLITE,
				Environment: "environments/prod",
				Activation:  true,
				DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: prodInstanceDir, Id: "admin"}},
			},
		})
		a.NoError(err)
		prodInstances = append(prodInstances, instance)
	}

	// Create issues that create databases.
	databaseName := "testTenantSchemaUpdate"
	for _, testInstance := range testInstances {
		err := ctl.createDatabaseV2(ctx, project, testInstance, nil, databaseName, "")
		a.NoError(err)
	}
	for _, prodInstance := range prodInstances {
		err := ctl.createDatabaseV2(ctx, project, prodInstance, nil, databaseName, "")
		a.NoError(err)
	}

	resp, err := ctl.databaseServiceClient.ListDatabases(ctx, &v1pb.ListDatabasesRequest{
		Parent: project.Name,
	})
	a.NoError(err)
	databases := resp.Databases

	var testDatabases []*v1pb.Database
	var prodDatabases []*v1pb.Database
	for _, testInstance := range testInstances {
		for _, database := range databases {
			if strings.HasPrefix(database.Name, testInstance.Name) {
				testDatabases = append(testDatabases, database)
				break
			}
		}
	}
	for _, prodInstance := range prodInstances {
		for _, database := range databases {
			if strings.HasPrefix(database.Name, prodInstance.Name) {
				prodDatabases = append(prodDatabases, database)
				break
			}
		}
	}
	a.Equal(testTenantNumber, len(testDatabases))
	a.Equal(prodTenantNumber, len(prodDatabases))

	databaseGroup, err := ctl.databaseGroupServiceClient.CreateDatabaseGroup(ctx, &v1pb.CreateDatabaseGroupRequest{
		Parent:          project.Name,
		DatabaseGroupId: "all",
		DatabaseGroup: &v1pb.DatabaseGroup{
			Title:        "all",
			DatabaseExpr: &expr.Expr{Expression: "true"},
		},
	})
	a.NoError(err)

	sheet, err := ctl.sheetServiceClient.CreateSheet(ctx, &v1pb.CreateSheetRequest{
		Parent: project.Name,
		Sheet: &v1pb.Sheet{
			Title:   "migration statement sheet",
			Content: []byte(migrationStatement1),
		},
	})
	a.NoError(err)

	// Create an issue that updates database schema.
	spec := &v1pb.Plan_Spec{
		Id: uuid.NewString(),
		Config: &v1pb.Plan_Spec_ChangeDatabaseConfig{
			ChangeDatabaseConfig: &v1pb.Plan_ChangeDatabaseConfig{
				Targets: []string{databaseGroup.Name},
				Sheet:   sheet.Name,
				Type:    v1pb.Plan_ChangeDatabaseConfig_MIGRATE,
			},
		},
	}
	plan, rollout, issue, err := ctl.changeDatabaseWithConfig(ctx, project, spec)
	a.NoError(err)

	// Assert the plan deployment.
	a.Len(plan.Deployment.DatabaseGroupMappings, 1)
	a.Equal(databaseGroup.Name, plan.Deployment.DatabaseGroupMappings[0].DatabaseGroup)
	a.Len(plan.Deployment.DatabaseGroupMappings[0].Databases, 2)
	a.ElementsMatch([]string{testDatabases[0].Name, prodDatabases[0].Name}, plan.Deployment.DatabaseGroupMappings[0].Databases)

	// Query schema.
	for _, testInstance := range testInstances {
		dbMetadata, err := ctl.databaseServiceClient.GetDatabaseSchema(ctx, &v1pb.GetDatabaseSchemaRequest{Name: fmt.Sprintf("%s/databases/%s/schema", testInstance.Name, databaseName)})
		a.NoError(err)
		a.Equal(wantBookSchema, dbMetadata.Schema)
	}
	for _, prodInstance := range prodInstances {
		dbMetadata, err := ctl.databaseServiceClient.GetDatabaseSchema(ctx, &v1pb.GetDatabaseSchemaRequest{Name: fmt.Sprintf("%s/databases/%s/schema", prodInstance.Name, databaseName)})
		a.NoError(err)
		a.Equal(wantBookSchema, dbMetadata.Schema)
	}

	// Create another database in the prod environment.
	databaseName2 := "testTenantSchemaUpdate2"
	err = ctl.createDatabaseV2(ctx, project, prodInstances[0], nil, databaseName2, "")
	a.NoError(err)

	resp, err = ctl.databaseServiceClient.ListDatabases(ctx, &v1pb.ListDatabasesRequest{
		Parent: project.Name,
	})
	a.NoError(err)
	databases = resp.Databases
	prodDatabases = nil
	for _, prodInstance := range prodInstances {
		for _, database := range databases {
			if strings.HasPrefix(database.Name, prodInstance.Name) {
				prodDatabases = append(prodDatabases, database)
			}
		}
	}
	a.Len(prodDatabases, 2)

	// Update the plan deployment.
	plan, err = ctl.planServiceClient.UpdatePlan(ctx, &v1pb.UpdatePlanRequest{
		Plan: &v1pb.Plan{
			Name: plan.Name,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"deployment"},
		},
	})
	a.NoError(err)

	a.Len(plan.Deployment.DatabaseGroupMappings, 1)
	a.Equal(databaseGroup.Name, plan.Deployment.DatabaseGroupMappings[0].DatabaseGroup)
	a.Len(plan.Deployment.DatabaseGroupMappings[0].Databases, 3)
	a.ElementsMatch([]string{testDatabases[0].Name, prodDatabases[0].Name, prodDatabases[1].Name}, plan.Deployment.DatabaseGroupMappings[0].Databases)

	// Create the new task.
	rollout2, err := ctl.rolloutServiceClient.CreateRollout(ctx, &v1pb.CreateRolloutRequest{
		Parent: project.Name,
		Rollout: &v1pb.Rollout{
			Plan: plan.Name,
		},
		Target: nil, // set to nil to create all stages and tasks.
	})
	a.NoError(err)
	a.Equal(rollout.Name, rollout2.Name)

	a.Len(rollout.Stages, 2)
	a.Len(rollout2.Stages, 2)
	a.Len(rollout.Stages[1].Tasks, 1)
	a.Len(rollout2.Stages[1].Tasks, 2)
	// The task for databaseName2 should be created.
	a.Equal(rollout2.Stages[1].Tasks[1].Target, fmt.Sprintf("%s/databases/%s", prodInstances[0].Name, databaseName2))

	err = ctl.waitRollout(ctx, issue.Name, rollout.Name)
	a.NoError(err)
}
