package database

import (
	"log/slog"

	"github.com/cloudstorex/backend/internal/identity"
	"github.com/cloudstorex/backend/internal/workspace"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	slog.Info("Running database migrations...")
	
	err := db.AutoMigrate(
		&identity.User{},
		&workspace.Workspace{},
	)
	
	if err != nil {
		return err
	}
	
	slog.Info("Database migrations completed successfully")
	return nil
}
