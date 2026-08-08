package serviceaccounts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceAccount struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	WorkspaceID    *uuid.UUID     `gorm:"type:uuid;index"` // Optional isolation
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    string         `gorm:"type:text"`
	ClientID       string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	ClientSecret   string         `gorm:"type:varchar(255);not null"` // Hashed
	Active         bool           `gorm:"default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
