package workspace

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkspaceType string

const (
	WorkspaceTypePersonal     WorkspaceType = "personal"
	WorkspaceTypeDeveloper    WorkspaceType = "developer"
	WorkspaceTypeOrganization WorkspaceType = "organization"
)

type Workspace struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID *uuid.UUID     `gorm:"type:uuid;index"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Type           WorkspaceType  `gorm:"type:varchar(50);not null"`
	OwnerID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
