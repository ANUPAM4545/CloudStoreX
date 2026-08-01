package database

import (
	"log/slog"

	"github.com/cloudstorex/backend/internal/identity"
	analyticsModel "github.com/cloudstorex/backend/internal/analytics/model"
	auditModel "github.com/cloudstorex/backend/internal/audit/model"
	jobsModel "github.com/cloudstorex/backend/internal/jobs/model"
	metadataModel "github.com/cloudstorex/backend/internal/metadata/model"
	policyModel "github.com/cloudstorex/backend/internal/policy/model"
	providerModel "github.com/cloudstorex/backend/internal/provider/model"
	quotaModel "github.com/cloudstorex/backend/internal/quota/model"
	"github.com/cloudstorex/backend/internal/workspace"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	slog.Info("Running database migrations...")
	
	err := db.AutoMigrate(
		&identity.User{},
		&workspace.Workspace{},
		&metadataModel.Bucket{},
		&metadataModel.Object{},
		&metadataModel.ObjectVersion{},
		&metadataModel.ObjectTag{},
		&metadataModel.ObjectMetadata{},
		&providerModel.Provider{},
		&policyModel.Policy{},
		&policyModel.RoutingDecision{},
		&jobsModel.Job{},
		&jobsModel.JobLog{},
		&quotaModel.WorkspaceQuota{},
		&auditModel.AuditLog{},
		&analyticsModel.StorageSnapshot{},
	)
	
	if err != nil {
		return err
	}
	
	slog.Info("Database migrations completed successfully")
	return nil
}
