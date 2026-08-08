package rbac

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;index"` // Null for system roles
	Name           string         `gorm:"type:varchar(255);not null"`
	Description    string         `gorm:"type:text"`
	Permissions    []string       `gorm:"type:jsonb"` // e.g. ["storage:bucket:read", "storage:object:*"]
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type RoleBinding struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoleID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	SubjectID uuid.UUID      `gorm:"type:uuid;not null;index"` // Can be UserID, UserGroupID, ServiceAccountID
	Scope     string         `gorm:"type:varchar(100);not null"` // e.g. "organization", "workspace", "bucket"
	ScopeID   uuid.UUID      `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
