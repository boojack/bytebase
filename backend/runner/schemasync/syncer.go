// Package schemasync is a runner that synchronize database schemas.
package schemasync

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sourcegraph/conc/pool"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/component/config"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	"github.com/bytebase/bytebase/backend/component/state"
	"github.com/bytebase/bytebase/backend/enterprise"
	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	v1pb "github.com/bytebase/bytebase/backend/generated-go/v1"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/store/model"
)

const (
	instanceSyncInterval        = 15 * time.Minute
	databaseSyncCheckerInterval = 10 * time.Second
	syncTimeout                 = 15 * time.Minute
	// defaultSyncInterval means never sync.
	defaultSyncInterval = 0 * time.Second
	MaximumOutstanding  = 100
)

// NewSyncer creates a schema syncer.
func NewSyncer(stores *store.Store, dbFactory *dbfactory.DBFactory, profile *config.Profile, stateCfg *state.State, licenseService *enterprise.LicenseService) *Syncer {
	return &Syncer{
		store:          stores,
		dbFactory:      dbFactory,
		profile:        profile,
		stateCfg:       stateCfg,
		licenseService: licenseService,
	}
}

// Syncer is the schema syncer.
type Syncer struct {
	sync.Mutex

	store           *store.Store
	dbFactory       *dbfactory.DBFactory
	profile         *config.Profile
	stateCfg        *state.State
	licenseService  *enterprise.LicenseService
	databaseSyncMap sync.Map // map[string]*store.DatabaseMessage
}

// Run will run the schema syncer once.
func (s *Syncer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	sp := pool.New()
	sp.Go(func() {
		slog.Debug(fmt.Sprintf("Schema syncer started and will run every %v", instanceSyncInterval))
		ticker := time.NewTicker(instanceSyncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.trySyncAll(ctx)
			case <-ctx.Done(): // if cancel() execute
				return
			}
		}
	})

	sp.Go(func() {
		ticker := time.NewTicker(databaseSyncCheckerInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				instances, err := s.store.ListInstancesV2(ctx, &store.FindInstanceMessage{})
				if err != nil {
					if err != nil {
						slog.Error("Failed to list instance", log.BBError(err))
						return
					}
				}
				instanceMap := make(map[string]*store.InstanceMessage)
				for _, instance := range instances {
					instanceMap[instance.ResourceID] = instance
				}
				dbwp := pool.New().WithMaxGoroutines(MaximumOutstanding)
				s.databaseSyncMap.Range(func(key, value any) bool {
					database, ok := value.(*store.DatabaseMessage)
					if !ok {
						return true
					}

					instance, ok := instanceMap[database.InstanceID]
					if !ok {
						slog.Debug("Instance not found",
							slog.String("instance", database.InstanceID),
							log.BBError(err))
						return true
					}
					maximumConnections := int(instance.Metadata.GetMaximumConnections())
					if maximumConnections <= 0 {
						maximumConnections = common.DefaultInstanceMaximumConnections
					}
					if s.stateCfg.InstanceOutstandingConnections.Increment(instance.ResourceID, maximumConnections) {
						return true
					}

					s.databaseSyncMap.Delete(key)
					dbwp.Go(func() {
						defer func() {
							s.stateCfg.InstanceOutstandingConnections.Decrement(instance.ResourceID)
						}()
						slog.Debug("Sync database schema", slog.String("instance", database.InstanceID), slog.String("database", database.DatabaseName))
						if err := s.SyncDatabaseSchema(ctx, database); err != nil {
							slog.Debug("Failed to sync database schema",
								slog.String("instance", database.InstanceID),
								slog.String("databaseName", database.DatabaseName),
								log.BBError(err))
						}
					})
					return true
				})
				dbwp.Wait()
			case <-ctx.Done(): // if cancel() execute
				return
			}
		}
	})
	sp.Wait()
}

func (s *Syncer) trySyncAll(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.Errorf("%v", r)
			}
			slog.Error("Instance syncer PANIC RECOVER", log.BBError(err), log.BBStack("panic-stack"))
		}
	}()

	wp := pool.New().WithMaxGoroutines(MaximumOutstanding)
	instances, err := s.store.ListInstancesV2(ctx, &store.FindInstanceMessage{})
	if err != nil {
		slog.Error("Failed to retrieve instances", log.BBError(err))
		return
	}
	now := time.Now()
	for _, instance := range instances {
		instance := instance
		interval := getOrDefaultSyncInterval(instance)
		if interval == defaultSyncInterval {
			continue
		}
		lastSyncTime := getOrDefaultLastSyncTime(instance.Metadata.LastSyncTime)
		// lastSyncTime + syncInterval > now
		// Next round not started yet.
		nextSyncTime := lastSyncTime.Add(interval)
		if now.Before(nextSyncTime) {
			continue
		}

		wp.Go(func() {
			slog.Debug("Sync instance schema", slog.String("instance", instance.ResourceID))
			if _, _, _, err := s.SyncInstance(ctx, instance); err != nil {
				slog.Debug("Failed to sync instance",
					slog.String("instance", instance.ResourceID),
					slog.String("error", err.Error()))
			}
		})
	}
	wp.Wait()

	instancesMap := map[string]*store.InstanceMessage{}
	for _, instance := range instances {
		instancesMap[instance.ResourceID] = instance
	}

	databases, err := s.store.ListDatabases(ctx, &store.FindDatabaseMessage{})
	if err != nil {
		slog.Error("Failed to retrieve databases", log.BBError(err))
		return
	}
	for _, database := range databases {
		database := database
		if database.Deleted {
			continue
		}
		instance, ok := instancesMap[database.InstanceID]
		if !ok {
			continue
		}
		// The database inherits the sync interval from the instance.
		interval := getOrDefaultSyncInterval(instance)
		if interval == defaultSyncInterval {
			continue
		}
		lastSyncTime := getOrDefaultLastSyncTime(database.Metadata.LastSyncTime)
		// lastSyncTime + syncInterval > now
		// Next round not started yet.
		nextSyncTime := lastSyncTime.Add(interval)
		if now.Before(nextSyncTime) {
			continue
		}

		s.databaseSyncMap.Store(database.String(), database)
	}
}

func (s *Syncer) SyncAllDatabases(ctx context.Context, instance *store.InstanceMessage) {
	find := &store.FindDatabaseMessage{}
	if instance != nil {
		find.InstanceID = &instance.ResourceID
	}
	databases, err := s.store.ListDatabases(ctx, find)
	if err != nil {
		slog.Debug("Failed to find databases to sync",
			slog.String("error", err.Error()))
		return
	}

	for _, database := range databases {
		// Skip deleted databases.
		if database.Deleted {
			continue
		}
		s.databaseSyncMap.Store(database.String(), database)
	}
}

func (s *Syncer) SyncDatabaseAsync(database *store.DatabaseMessage) {
	if database == nil || database.Deleted {
		return
	}
	s.databaseSyncMap.Store(database.String(), database)
}

func (s *Syncer) SyncDatabasesAsync(databases []*store.DatabaseMessage) {
	for _, database := range databases {
		s.SyncDatabaseAsync(database)
	}
}

// GetInstanceMeta gets the instance metadata.
func (s *Syncer) GetInstanceMeta(ctx context.Context, instance *store.InstanceMessage) (*db.InstanceMetadata, error) {
	driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, instance, nil /* database */, db.ConnectionContext{})
	if err != nil {
		return nil, err
	}
	defer driver.Close(ctx)

	deadlineCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(syncTimeout))
	defer cancelFunc()
	instanceMeta, err := driver.SyncInstance(deadlineCtx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to sync instance: %s", instance.ResourceID)
	}

	if instanceMeta.Metadata == nil {
		instanceMeta.Metadata = &storepb.Instance{}
	}

	instanceMeta.Metadata.LastSyncTime = timestamppb.Now()

	return instanceMeta, nil
}

// SyncInstance syncs the schema for all databases in an instance.
func (s *Syncer) SyncInstance(ctx context.Context, instance *store.InstanceMessage) (*store.InstanceMessage, []*storepb.DatabaseSchemaMetadata, []*store.DatabaseMessage, error) {
	instanceMeta, err := s.GetInstanceMeta(ctx, instance)
	if err != nil {
		return nil, nil, nil, err
	}
	metadata, ok := proto.Clone(instance.Metadata).(*storepb.Instance)
	if !ok {
		return nil, nil, nil, errors.Errorf("failed to convert instance metadata type")
	}
	metadata.LastSyncTime = instanceMeta.Metadata.LastSyncTime
	metadata.MysqlLowerCaseTableNames = instanceMeta.Metadata.MysqlLowerCaseTableNames
	metadata.Roles = instanceMeta.Metadata.Roles

	updateInstance := &store.UpdateInstanceMessage{
		ResourceID: instance.ResourceID,
		Metadata:   metadata,
	}
	if instanceMeta.Version != instance.Metadata.GetVersion() {
		metadata.Version = instanceMeta.Version
	}
	updatedInstance, err := s.store.UpdateInstanceV2(ctx, updateInstance)
	if err != nil {
		return nil, nil, nil, err
	}

	databases, err := s.store.ListDatabases(ctx, &store.FindDatabaseMessage{InstanceID: &instance.ResourceID})
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to sync database for instance: %s. Failed to find database list", instance.ResourceID)
	}
	var newDatabases []*store.DatabaseMessage
	var filteredDatabaseMetadatas []*storepb.DatabaseSchemaMetadata

	for _, databaseMetadata := range instanceMeta.Databases {
		if len(instance.Metadata.GetSyncDatabases()) > 0 && !slices.Contains(instance.Metadata.GetSyncDatabases(), databaseMetadata.Name) {
			continue
		}
		filteredDatabaseMetadatas = append(filteredDatabaseMetadatas, databaseMetadata)
		idx := slices.IndexFunc(databases, func(db *store.DatabaseMessage) bool { return db.DatabaseName == databaseMetadata.Name })

		if idx < 0 {
			// Create the database in the default project.
			newDatabase, err := s.store.CreateDatabaseDefault(ctx, &store.DatabaseMessage{
				InstanceID:   instance.ResourceID,
				DatabaseName: databaseMetadata.Name,
				ProjectID:    common.DefaultProjectID,
			})
			if err != nil {
				return nil, nil, nil, errors.Wrapf(err, "failed to create instance %q database %q in sync runner", instance.ResourceID, databaseMetadata.Name)
			}
			if newDatabase != nil {
				newDatabases = append(newDatabases, newDatabase)
			}
		}
	}

	for _, database := range databases {
		idx := slices.IndexFunc(filteredDatabaseMetadatas, func(db *storepb.DatabaseSchemaMetadata) bool { return db.Name == database.DatabaseName })
		if idx < 0 {
			d := true
			if _, err := s.store.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
				InstanceID:   instance.ResourceID,
				DatabaseName: database.DatabaseName,
				Deleted:      &d,
			}); err != nil {
				return nil, nil, nil, errors.Errorf("failed to update database %q for instance %q", database.DatabaseName, instance.ResourceID)
			}
		}
	}

	return updatedInstance, instanceMeta.Databases, newDatabases, nil
}

// SyncDatabaseSchema will sync the schema for a database.
func (s *Syncer) SyncDatabaseSchemaToHistory(ctx context.Context, database *store.DatabaseMessage) (int64, error) {
	instance, err := s.store.GetInstanceV2(ctx, &store.FindInstanceMessage{ResourceID: &database.InstanceID})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get instance %q", database.InstanceID)
	}
	if instance == nil {
		return 0, errors.Errorf("instance %q not found", database.InstanceID)
	}
	driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, instance, database, db.ConnectionContext{})
	if err != nil {
		return 0, err
	}
	defer driver.Close(ctx)
	// Sync database schema
	deadlineCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(syncTimeout))
	defer cancelFunc()
	databaseMetadata, err := driver.SyncDBSchema(deadlineCtx)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to sync database schema for database %q", database.DatabaseName)
	}

	dbSchema, err := s.store.GetDBSchema(ctx, database.InstanceID, database.DatabaseName)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get database schema for database %q", database.DatabaseName)
	}

	dbModelConfig := model.NewDatabaseConfig(nil)
	if dbSchema != nil {
		dbModelConfig = dbSchema.GetInternalConfig()
	}

	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &database.ProjectID,
	})
	if err != nil {
		return 0, errors.Wrapf(err, `failed to get project by id "%s"`, database.ProjectID)
	}
	classificationConfig, err := s.store.GetDataClassificationConfigByID(ctx, project.DataClassificationConfigID)
	if err != nil {
		return 0, errors.Wrapf(err, `failed to get classification config by id "%s"`, project.DataClassificationConfigID)
	}

	if instance.Metadata.GetEngine() != storepb.Engine_MYSQL && instance.Metadata.GetEngine() != storepb.Engine_POSTGRES {
		// Force to disable classification from comment if the engine is not MYSQL or PG.
		classificationConfig.ClassificationFromConfig = true
	}
	var dbConfig *storepb.DatabaseConfig
	if classificationConfig.ClassificationFromConfig {
		// Only set the user comment.
		setUserCommentFromComment(databaseMetadata)
		dbConfig = dbModelConfig.BuildDatabaseConfig()
	} else {
		// Get classification from the comment.
		if err := s.licenseService.IsFeatureEnabledForInstance(v1pb.PlanFeature_FEATURE_DATA_CLASSIFICATION, instance); err == nil {
			dbConfig = buildDatabaseConfigWithClassificationFromComment(databaseMetadata, dbModelConfig, classificationConfig)
		} else {
			dbConfig = dbModelConfig.BuildDatabaseConfig()
		}
	}

	d := false
	metadata, ok := proto.Clone(database.Metadata).(*storepb.DatabaseMetadata)
	if !ok {
		return 0, errors.Errorf("failed to convert database metadata type")
	}
	metadata.LastSyncTime = timestamppb.New(time.Now())
	metadata.BackupAvailable = s.databaseBackupAvailable(ctx, instance, databaseMetadata)
	metadata.Datashare = databaseMetadata.Datashare
	if _, err := s.store.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
		InstanceID:   database.InstanceID,
		DatabaseName: database.DatabaseName,
		Deleted:      &d,
		Metadata:     metadata,
	}); err != nil {
		return 0, errors.Wrapf(err, "failed to update database %q for instance %q", database.DatabaseName, database.InstanceID)
	}
	var schemaBuf bytes.Buffer
	if err := driver.Dump(ctx, &schemaBuf, databaseMetadata); err != nil {
		return 0, errors.Wrapf(err, "failed to dump database schema for database %q", database.DatabaseName)
	}
	rawDump := schemaBuf.Bytes()

	todo := !common.EngineDBSchemaReadyToMigrate(instance.Metadata.GetEngine())
	if err := s.store.UpsertDBSchema(ctx,
		database.InstanceID, database.DatabaseName,
		databaseMetadata, dbConfig, rawDump, todo,
	); err != nil {
		if strings.Contains(err.Error(), "escape sequence") {
			if metadataBytes, err := protojson.Marshal(databaseMetadata); err == nil {
				slog.Error("unsupported Unicode escape sequence", slog.String("metadata", string(metadataBytes)), slog.String("raw_dump", string(rawDump)))
			}
		}
		return 0, errors.Wrapf(err, "failed to upsert database schema for database %q", database.DatabaseName)
	}

	id, err := s.store.CreateSyncHistory(ctx, database.InstanceID, database.DatabaseName, databaseMetadata, string(rawDump))
	if err != nil {
		if strings.Contains(err.Error(), "escape sequence") {
			if metadataBytes, err := protojson.Marshal(databaseMetadata); err == nil {
				slog.Error("unsupported Unicode escape sequence", slog.String("metadata", string(metadataBytes)), slog.String("raw_dump", string(rawDump)))
			}
		}
		return 0, errors.Wrapf(err, "failed to insert sync history for database %q", database.DatabaseName)
	}

	return id, nil
}

// SyncDatabaseSchema will sync the schema for a database.
func (s *Syncer) SyncDatabaseSchema(ctx context.Context, database *store.DatabaseMessage) (retErr error) {
	instance, err := s.store.GetInstanceV2(ctx, &store.FindInstanceMessage{ResourceID: &database.InstanceID})
	if err != nil {
		return errors.Wrapf(err, "failed to get instance %q", database.InstanceID)
	}
	if instance == nil {
		return errors.Errorf("instance %q not found", database.InstanceID)
	}
	driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, instance, database, db.ConnectionContext{})
	if err != nil {
		return err
	}
	defer driver.Close(ctx)
	// Sync database schema
	deadlineCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(syncTimeout))
	defer cancelFunc()
	databaseMetadata, err := driver.SyncDBSchema(deadlineCtx)
	if err != nil {
		return errors.Wrapf(err, "failed to sync database schema for database %q", database.DatabaseName)
	}
	var schemaBuf bytes.Buffer
	if err := driver.Dump(deadlineCtx, &schemaBuf, databaseMetadata); err != nil {
		return errors.Wrapf(err, "failed to dump database schema for database %q", database.DatabaseName)
	}
	rawDump := schemaBuf.Bytes()

	dbSchema, err := s.store.GetDBSchema(ctx, database.InstanceID, database.DatabaseName)
	if err != nil {
		return errors.Wrapf(err, "failed to get database schema for database %q", database.DatabaseName)
	}

	dbModelConfig := model.NewDatabaseConfig(nil)
	if dbSchema != nil {
		dbModelConfig = dbSchema.GetInternalConfig()
	}

	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &database.ProjectID,
	})
	if err != nil {
		return errors.Wrapf(err, `failed to get project by id "%s"`, database.ProjectID)
	}
	classificationConfig, err := s.store.GetDataClassificationConfigByID(ctx, project.DataClassificationConfigID)
	if err != nil {
		return errors.Wrapf(err, `failed to get classification config by id "%s"`, project.DataClassificationConfigID)
	}

	if instance.Metadata.GetEngine() != storepb.Engine_MYSQL && instance.Metadata.GetEngine() != storepb.Engine_POSTGRES {
		// Force to disable classification from comment if the engine is not MYSQL or PG.
		classificationConfig.ClassificationFromConfig = true
	}
	var dbConfig *storepb.DatabaseConfig
	if classificationConfig.ClassificationFromConfig {
		// Only set the user comment.
		setUserCommentFromComment(databaseMetadata)
		dbConfig = dbModelConfig.BuildDatabaseConfig()
	} else {
		// Get classification from the comment.
		dbConfig = buildDatabaseConfigWithClassificationFromComment(databaseMetadata, dbModelConfig, classificationConfig)
	}

	d := false
	metadata, ok := proto.Clone(database.Metadata).(*storepb.DatabaseMetadata)
	if !ok {
		return errors.Errorf("failed to convert database metadata type")
	}
	metadata.LastSyncTime = timestamppb.New(time.Now())
	metadata.BackupAvailable = s.databaseBackupAvailable(ctx, instance, databaseMetadata)
	metadata.Datashare = databaseMetadata.Datashare
	drifted, skipped, err := s.getSchemaDrifted(ctx, instance, database, string(rawDump))
	if err != nil {
		return errors.Wrapf(err, "failed to get schema drifted for database %q", database.DatabaseName)
	}
	if !skipped {
		metadata.Drifted = drifted
	}
	if _, err := s.store.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
		InstanceID:   database.InstanceID,
		DatabaseName: database.DatabaseName,
		Deleted:      &d,
		Metadata:     metadata,
	}); err != nil {
		return errors.Wrapf(err, "failed to update database %q for instance %q", database.DatabaseName, database.InstanceID)
	}

	todo := !common.EngineDBSchemaReadyToMigrate(instance.Metadata.GetEngine())
	if err := s.store.UpsertDBSchema(ctx,
		database.InstanceID, database.DatabaseName,
		databaseMetadata, dbConfig, rawDump, todo,
	); err != nil {
		if strings.Contains(err.Error(), "escape sequence") {
			if metadataBytes, err := protojson.Marshal(databaseMetadata); err == nil {
				slog.Error("unsupported Unicode escape sequence", slog.String("metadata", string(metadataBytes)), slog.String("raw_dump", string(rawDump)))
			}
		}
		return errors.Wrapf(err, "failed to upsert database schema for database %q", database.DatabaseName)
	}

	return nil
}

func (s *Syncer) getSchemaDrifted(ctx context.Context, instance *store.InstanceMessage, database *store.DatabaseMessage, rawDump string) (drifted bool, skipped bool, err error) {
	// Redis and MongoDB are schemaless.
	if disableSchemaDriftCheck(instance.Metadata.GetEngine()) {
		return false, false, nil
	}
	limit := 1
	list, err := s.store.ListChangelogs(ctx, &store.FindChangelogMessage{
		InstanceID:     &database.InstanceID,
		DatabaseName:   &database.DatabaseName,
		TypeList:       []string{string(db.Migrate), string(db.Baseline)},
		HasSyncHistory: true,
		Limit:          &limit,
		ShowFull:       true,
	})
	if err != nil {
		return false, false, errors.Wrapf(err, "failed to list changelogs")
	}
	if len(list) == 0 {
		return false, false, nil
	}

	changelog := list[0]
	if changelog.Payload.GetGitCommit() != s.profile.GitCommit {
		return false, true, nil
	}
	if changelog.SyncHistoryUID == nil {
		return false, false, errors.Errorf("expect sync history but get nil")
	}
	latestSchema := string(rawDump)
	return changelog.Schema != latestSchema, false, nil
}

func (s *Syncer) databaseBackupAvailable(ctx context.Context, instance *store.InstanceMessage, dbSchema *storepb.DatabaseSchemaMetadata) bool {
	if !common.EngineSupportPriorBackup(instance.Metadata.GetEngine()) {
		return false
	}
	switch instance.Metadata.GetEngine() {
	case storepb.Engine_POSTGRES:
		if dbSchema == nil {
			return false
		}
		for _, schema := range dbSchema.Schemas {
			if schema.GetName() == common.BackupDatabaseNameOfEngine(storepb.Engine_POSTGRES) {
				return true
			}
		}
	case storepb.Engine_MYSQL, storepb.Engine_MSSQL, storepb.Engine_TIDB:
		dbName := common.BackupDatabaseNameOfEngine(instance.Metadata.GetEngine())
		backupDB, err := s.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
			InstanceID:   &instance.ResourceID,
			DatabaseName: &dbName,
		})
		if err != nil {
			slog.Debug("Failed to get backup database", "err", err)
			return false
		}
		return backupDB != nil
	case storepb.Engine_ORACLE:
		dbName := common.BackupDatabaseNameOfEngine(storepb.Engine_ORACLE)
		backupDB, err := s.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
			InstanceID:   &instance.ResourceID,
			DatabaseName: &dbName,
		})
		if err != nil {
			slog.Debug("Failed to get backup database", "err", err)
			return false
		}
		return backupDB != nil
	default:
		// Unsupported database engine for backup
		slog.Debug("Unsupported database engine for backup", "engine", instance.Metadata.GetEngine())
		return false
	}
	return false
}

func buildDatabaseConfigWithClassificationFromComment(dbSchema *storepb.DatabaseSchemaMetadata, databaseConfig *model.DatabaseConfig, classificationConfig *storepb.DataClassificationSetting_DataClassificationConfig) *storepb.DatabaseConfig {
	// Create a new DatabaseConfig to build from scratch
	newConfig := &storepb.DatabaseConfig{
		Name: dbSchema.Name,
	}

	for _, schema := range dbSchema.Schemas {
		var tables []*storepb.TableCatalog
		hasSchemaContent := false

		// Get existing schema config to preserve SemanticType and Labels
		schemaConfig := databaseConfig.GetSchemaConfig(schema.Name)

		for _, table := range schema.Tables {
			classification, userComment := common.GetClassificationAndUserComment(table.Comment, classificationConfig)

			table.UserComment = userComment

			// Get existing table config
			tableConfig := schemaConfig.GetTableConfig(table.Name)

			var columns []*storepb.ColumnCatalog
			hasTableContent := false

			// Check if table has classification
			if classification != "" {
				hasTableContent = true
			}

			for _, col := range table.Columns {
				colClassification, colUserComment := common.GetClassificationAndUserComment(col.Comment, classificationConfig)

				col.UserComment = colUserComment

				// Get existing column config to preserve SemanticType and Labels
				existingColumnConfig := tableConfig.GetColumnConfig(col.Name)

				// Add column config if it has any meaningful data
				if colClassification != "" || existingColumnConfig.SemanticType != "" || len(existingColumnConfig.Labels) > 0 {
					columns = append(columns, &storepb.ColumnCatalog{
						Name:           col.Name,
						Classification: colClassification,
						SemanticType:   existingColumnConfig.SemanticType,
						Labels:         existingColumnConfig.Labels,
					})
					hasTableContent = true
				}
			}

			// Only add table catalog if it has content (either table classification or column configurations)
			if hasTableContent {
				tableCatalog := &storepb.TableCatalog{
					Name:           table.Name,
					Classification: classification,
					Columns:        columns,
				}
				tables = append(tables, tableCatalog)
				hasSchemaContent = true
			}
		}

		// Only add schema catalog if it has content
		if hasSchemaContent {
			schemaCatalog := &storepb.SchemaCatalog{
				Name:   schema.Name,
				Tables: tables,
			}
			newConfig.Schemas = append(newConfig.Schemas, schemaCatalog)
		}
	}

	return newConfig
}

func setUserCommentFromComment(dbSchema *storepb.DatabaseSchemaMetadata) {
	for _, schema := range dbSchema.Schemas {
		for _, table := range schema.Tables {
			table.UserComment = table.Comment
			for _, col := range table.Columns {
				col.UserComment = col.Comment
			}
		}
	}
}

func getOrDefaultSyncInterval(instance *store.InstanceMessage) time.Duration {
	if !instance.Metadata.GetActivation() {
		return defaultSyncInterval
	}
	if !instance.Metadata.GetSyncInterval().IsValid() {
		return defaultSyncInterval
	}
	if instance.Metadata.GetSyncInterval().GetSeconds() == 0 && instance.Metadata.GetSyncInterval().GetNanos() == 0 {
		return defaultSyncInterval
	}
	return instance.Metadata.GetSyncInterval().AsDuration()
}

func getOrDefaultLastSyncTime(t *timestamppb.Timestamp) time.Time {
	if t.IsValid() {
		return t.AsTime()
	}
	return time.Unix(0, 0)
}

func disableSchemaDriftCheck(dbTp storepb.Engine) bool {
	m := map[storepb.Engine]struct{}{
		storepb.Engine_MONGODB:    {},
		storepb.Engine_REDIS:      {},
		storepb.Engine_REDSHIFT:   {},
		storepb.Engine_RISINGWAVE: {},
	}
	_, ok := m[dbTp]
	return ok
}
