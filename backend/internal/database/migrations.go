package database

import (
	"log/slog"

	"github.com/cloudstorex/backend/internal/identity"
	metadataModel "github.com/cloudstorex/backend/internal/metadata/model"
	policyModel "github.com/cloudstorex/backend/internal/policy/model"
	providerModel "github.com/cloudstorex/backend/internal/provider/model"
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
		&metadataModel.ObjectTag{},
		&metadataModel.ObjectMetadata{},
		&providerModel.Provider{},
		&policyModel.Policy{},
		&policyModel.RoutingDecision{},
	)
	
	if err != nil {
		return err
	}
	
	slog.Info("Database migrations completed successfully")
	return nil
}
