package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Domain    string         `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Department struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name           string         `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Team struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	DepartmentID   *uuid.UUID     `gorm:"type:uuid;index"`
	Name           string         `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type UserGroup struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name           string         `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type EntityType string

const (
	EntityOrganization EntityType = "organization"
	EntityDepartment   EntityType = "department"
	EntityTeam         EntityType = "team"
	EntityGroup        EntityType = "group"
)

type Membership struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	EntityType EntityType     `gorm:"type:varchar(50);not null;index"`
	EntityID   uuid.UUID      `gorm:"type:uuid;not null;index"`
	Role       string         `gorm:"type:varchar(100);not null"` // e.g. "admin", "member"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type InvitationStatus string

const (
	InvitationPending  InvitationStatus = "pending"
	InvitationAccepted InvitationStatus = "accepted"
	InvitationExpired  InvitationStatus = "expired"
)

type Invitation struct {
	ID             uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID        `gorm:"type:uuid;not null;index"`
	Email          string           `gorm:"type:varchar(255);not null;index"`
	Role           string           `gorm:"type:varchar(100);not null"`
	Token          string           `gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt      time.Time        `gorm:"not null"`
	Status         InvitationStatus `gorm:"type:varchar(50);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
