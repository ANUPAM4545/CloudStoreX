package apikeys

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APIKey struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Prefix    string         `gorm:"type:varchar(20);not null;uniqueIndex:idx_prefix"` // To identify the key quickly
	KeyHash   string         `gorm:"type:varchar(255);not null"` // argon2 hash of the secret key
	Scopes    []string       `gorm:"type:jsonb"`
	ExpiresAt *time.Time
	LastUsed  *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
