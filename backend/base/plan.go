package base

import (
	"fmt"
)

// PlanType is the type for a plan.
type PlanType int

const (
	// FREE is the plan type for FREE.
	FREE PlanType = iota
	// TEAM is the plan type for TEAM.
	TEAM
	// ENTERPRISE is the plan type for ENTERPRISE.
	ENTERPRISE
)

// String returns the string format of plan type.
func (p PlanType) String() string {
	switch p {
	case FREE:
		return "FREE"
	case TEAM:
		return "TEAM"
	case ENTERPRISE:
		return "ENTERPRISE"
	}
	return ""
}

// Priority returns the priority of the plan type.
// Higher priority means the plan supports more features.
func (p PlanType) Priority() int {
	switch p {
	case FREE:
		return 1
	case TEAM:
		return 2
	case ENTERPRISE:
		return 3
	}
	return 0
}

// FeatureType is the type of a feature.
type FeatureType string

const (
	// Admin & Security.

	// FeatureSSO allows user to manage SSO provider and authenticate (login) with SSO.
	FeatureSSO FeatureType = "bb.feature.sso"
	// Feature2FA allows user to manage 2FA provider and authenticate (login) with 2FA.
	Feature2FA FeatureType = "bb.feature.2fa"
	// FeatureDisallowSignup allows user to change the disallow signup flag.
	FeatureDisallowSignup FeatureType = "bb.feature.disallow-signup"
	// FeatureDisallowPasswordSignin allows user to disallow password signin.
	FeatureDisallowPasswordSignin FeatureType = "bb.feature.disallow-password-signin"
	// FeatureSecureToken allows user to manage authentication token security.
	FeatureSecureToken FeatureType = "bb.feature.secure-token"
	// FeatureExternalSecretManager uses secrets from external secret manager.
	FeatureExternalSecretManager FeatureType = "bb.feature.external-secret-manager"
	// FeaturePasswordRestriction allows user to configure the password restriction.
	FeaturePasswordRestriction FeatureType = "bb.feature.password-restriction"
	// FeatureDomainRestriction allows user to enforce the member email with workspace domain.
	FeatureDomainRestriction FeatureType = "bb.feature.domain-restriction"
	// FeatureDirectorySync allows to sync users and groups from Entra ID.
	FeatureDirectorySync FeatureType = "bb.feature.directory-sync"

	// FeatureRBAC enables RBAC.
	//
	// - Workspace level RBAC
	// - Project level RBAC.
	FeatureRBAC FeatureType = "bb.feature.rbac"

	// FeatureWatermark enables full-screen watermark.
	FeatureWatermark FeatureType = "bb.feature.watermark"

	// FeatureAuditLog enables viewing audit logs.
	FeatureAuditLog FeatureType = "bb.feature.audit-log"

	// FeatureCustomRole enables customizing roles.
	FeatureCustomRole FeatureType = "bb.feature.custom-role"

	// FeatureIssueAdvancedSearch supports search issue with advanced filter.
	FeatureIssueAdvancedSearch FeatureType = "bb.feature.issue-advanced-search"

	// FeatureAnnouncement enable announcement banner setting.
	FeatureAnnouncement FeatureType = "bb.feature.announcement"

	// Branding.

	// FeatureBranding enables customized branding.
	//
	// Currently, we only support customizing the logo.
	FeatureBranding FeatureType = "bb.feature.branding"

	// Change Workflow.

	// FeatureDBAWorkflow enforces the DBA workflow.
	//
	// - Developers can't create and view instances since they are exclusively by DBA, they can
	//   only access database.
	// - Developers can't create database.
	// - Developers can't query and export data directly. They must request corresponding permissions first.
	FeatureDBAWorkflow FeatureType = "bb.feature.dba-workflow"
	// FeatureSchemaDrift detects if there occurs schema drift.
	// See https://bytebase.com/docs/features/drift-detection
	FeatureSchemaDrift FeatureType = "bb.feature.schema-drift"
	// FeatureSQLReview allows user to specify schema policy for the environment
	//
	// e.g. One can configure rules for database schema or SQL query.
	FeatureSQLReview FeatureType = "bb.feature.sql-review"
	// FeatureTaskScheduleTime allows user to run task at a scheduled time.
	// FeatureEncryptedSecrets is a feature that allows user to setting the encrypted secrets for the database.
	FeatureEncryptedSecrets FeatureType = "bb.feature.encrypted-secrets"
	// FeatureDatabaseGrouping allows user to create database/schema groups.
	FeatureDatabaseGrouping FeatureType = "bb.feature.database-grouping"
	// FeatureSchemaTemplate allows user to create and use the schema template.
	FeatureSchemaTemplate FeatureType = "bb.feature.schema-template"
	// FeatureIssueProjectSetting supports some advanced project settings for issue.
	FeatureIssueProjectSetting FeatureType = "bb.feature.issue-project-setting"

	// Database management.

	// FeatureReadReplicaConnection allows user to set a read replica connection
	// including host and port to data source.
	FeatureReadReplicaConnection FeatureType = "bb.feature.read-replica-connection"
	// FeatureCustomInstanceSynchronization allows user to customize the synchronization for instance.
	FeatureCustomInstanceSynchronization FeatureType = "bb.feature.custom-instance-synchronization"
	// FeatureSyncSchemaAllVersions allows user to sync the base database schema all versions into target database.
	FeatureSyncSchemaAllVersions FeatureType = "bb.feature.sync-schema-all-versions"

	// Policy Control.

	// FeatureRolloutPolicy allows user to specify approval policy for the environment
	//
	// e.g. One can configure to NOT require approval for dev environment while require
	//      manual rollout for production.
	FeatureRolloutPolicy FeatureType = "bb.feature.rollout-policy"
	// FeatureEnvironmentTierPolicy allows user to set the tier of an environment.
	//
	// e.g. set the tier to "PROTECTED" for the production environment.
	FeatureEnvironmentTierPolicy FeatureType = "bb.feature.environment-tier-policy"
	// FeatureSensitiveData allows user to annotate and protect sensitive data.
	FeatureSensitiveData FeatureType = "bb.feature.sensitive-data"
	// FeatureAccessControl allows user to config the access control.
	FeatureAccessControl FeatureType = "bb.feature.access-control"
	// FeatureCustomApproval enables custom risk level definition and custom
	// approval chain definition.
	FeatureCustomApproval FeatureType = "bb.feature.custom-approval"

	// Efficiency.

	// FeatureBatchQuery enables batch query databases in SQL Editor.
	FeatureBatchQuery FeatureType = "bb.feature.batch-query"

	// Collaboration.

	// FeatureSharedSQLScript enables sharing sql script.
	FeatureSharedSQLScript FeatureType = "bb.feature.shared-sql-script"

	// AI Assistant.

	// FeatureAIAssistant enables AI features powered by OpenAI.
	FeatureAIAssistant FeatureType = "bb.feature.ai-assistant"
)

func (e FeatureType) String() string {
	return string(e)
}

// Name returns a readable name of the feature.
func (e FeatureType) Name() string {
	switch e {
	// Admin & Security
	case FeatureSSO:
		return "SSO"
	case Feature2FA:
		return "2FA"
	case FeatureDisallowSignup:
		return "Disallow singup"
	case FeatureDisallowPasswordSignin:
		return "Disallow password signin"
	case FeatureSecureToken:
		return "Secure token"
	case FeaturePasswordRestriction:
		return "Password restriction"
	case FeatureDomainRestriction:
		return "Domain restriction"
	case FeatureDirectorySync:
		return "Directory sync"
	case FeatureExternalSecretManager:
		return "External Secret Manager"
	case FeatureRBAC:
		return "RBAC"
	case FeatureWatermark:
		return "Watermark"
	case FeatureAuditLog:
		return "Audit log"
	case FeatureCustomRole:
		return "Custom role"
	case FeatureIssueAdvancedSearch:
		return "Advanced search"
	case FeatureAnnouncement:
		return "Announcement"
	// Branding
	case FeatureBranding:
		return "Branding"
	// Change Workflow
	case FeatureDBAWorkflow:
		return "DBA workflow"
	case FeatureSchemaDrift:
		return "Schema drift"
	case FeatureSQLReview:
		return "SQL review"
	case FeatureEncryptedSecrets:
		return "Encrypted secrets"
	case FeatureDatabaseGrouping:
		return "Database grouping"
	case FeatureSchemaTemplate:
		return "Schema template"
	case FeatureIssueProjectSetting:
		return "Issue project setting"
	// Database management
	case FeatureReadReplicaConnection:
		return "Read replica connection"
	case FeatureCustomInstanceSynchronization:
		return "Custom synchronization for instance"
	case FeatureSyncSchemaAllVersions:
		return "Synchronize schema all versions"
	// Policy Control
	case FeatureRolloutPolicy:
		return "Rollout policy"
	case FeatureEnvironmentTierPolicy:
		return "Environment tier"
	case FeatureSensitiveData:
		return "Sensitive data"
	case FeatureAccessControl:
		return "Access Control"
	case FeatureCustomApproval:
		return "Custom Approval"
	// Efficiency
	case FeatureBatchQuery:
		return "Batch query"
	// Collaboration
	case FeatureSharedSQLScript:
		return "Shared SQL script"
	// Plugins
	case FeatureAIAssistant:
		return "OpenAI"
	}
	return ""
}

// AccessErrorMessage returns a error message with feature name and minimum supported plan.
func (e FeatureType) AccessErrorMessage() string {
	plan := e.minimumSupportedPlan()
	return fmt.Sprintf("%s is a %s feature, please upgrade to access it.", e.Name(), plan.String())
}

// minimumSupportedPlan will find the minimum plan which supports the target feature.
func (e FeatureType) minimumSupportedPlan() PlanType {
	for i, enabled := range FeatureMatrix[e] {
		if enabled {
			return PlanType(i)
		}
	}

	return ENTERPRISE
}

// FeatureMatrix is a map from the a particular feature to the respective enablement of a particular
// plan in [FREE, TEAM, Enterprise].
var FeatureMatrix = map[FeatureType][3]bool{
	// Admin & Security
	FeatureSSO:                    {false, true, true},
	Feature2FA:                    {false, false, true},
	FeatureDisallowSignup:         {false, true, true},
	FeatureDisallowPasswordSignin: {false, false, true},
	FeatureSecureToken:            {false, false, true},
	FeaturePasswordRestriction:    {false, false, true},
	FeatureDomainRestriction:      {false, false, true},
	FeatureDirectorySync:          {false, false, true},
	FeatureExternalSecretManager:  {false, false, true},
	FeatureRBAC:                   {true, true, true},
	FeatureWatermark:              {false, false, true},
	FeatureAuditLog:               {false, false, true},
	FeatureCustomRole:             {false, false, true},
	FeatureIssueAdvancedSearch:    {false, true, true},
	FeatureAnnouncement:           {false, false, true},
	// Branding
	FeatureBranding: {false, false, true},
	// Change Workflow
	FeatureDBAWorkflow:         {false, false, true},
	FeatureSchemaDrift:         {false, false, true},
	FeatureSQLReview:           {true, true, true},
	FeatureEncryptedSecrets:    {false, true, true},
	FeatureDatabaseGrouping:    {false, true, true},
	FeatureSchemaTemplate:      {false, false, true},
	FeatureIssueProjectSetting: {false, true, true},
	// Database management
	FeatureReadReplicaConnection:         {false, false, true},
	FeatureCustomInstanceSynchronization: {false, false, true},
	FeatureSyncSchemaAllVersions:         {false, true, true},
	// Policy Control
	FeatureRolloutPolicy:         {true, true, true},
	FeatureEnvironmentTierPolicy: {false, false, true},
	FeatureSensitiveData:         {false, false, true},
	FeatureAccessControl:         {false, false, true},
	FeatureCustomApproval:        {false, false, true},
	// Efficiency
	FeatureBatchQuery: {false, true, true},
	// Collaboration
	FeatureSharedSQLScript: {true, true, true},
	// Plugins
	FeatureAIAssistant: {true, true, true},
}

// InstanceLimitFeature is the map for instance feature. Only allowed to access these feature for activate instance.
var InstanceLimitFeature = map[FeatureType]bool{
	// Security
	FeatureExternalSecretManager: true,
	// Change Workflow
	FeatureSchemaDrift:      true,
	FeatureEncryptedSecrets: true,
	// Database management
	FeatureReadReplicaConnection:         true,
	FeatureCustomInstanceSynchronization: true,
	FeatureSyncSchemaAllVersions:         true,
	// Policy Control
	FeatureSensitiveData: true,
}

// Feature returns whether a particular feature is available in a particular plan.
func Feature(feature FeatureType, plan PlanType) bool {
	matrix, ok := FeatureMatrix[feature]
	if !ok {
		return false
	}
	return matrix[plan]
}
