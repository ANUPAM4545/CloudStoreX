package model

import (
	"time"

	"github.com/google/uuid"
)

// ObjectVersion represents a specific point-in-time version of an object.
type ObjectVersion struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ObjectID          uuid.UUID `gorm:"type:uuid;not null;index" json:"object_id"` // Reference to the parent Object
	VersionNumber     int       `gorm:"not null" json:"version_number"`
	ProviderID        string    `gorm:"type:varchar(255);not null" json:"provider_id"`
	ProviderObjectKey string    `gorm:"type:text;not null" json:"provider_object_key"`
	SizeBytes         int64     `gorm:"not null;default:0" json:"size_bytes"`
	ETag              string    `gorm:"type:varchar(255)" json:"etag"`
	IsCurrent         bool      `gorm:"not null;default:false;index" json:"is_current"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Retention & Compliance
	RetainUntil *time.Time `gorm:"index" json:"retain_until"`
	LegalHold   bool       `gorm:"not null;default:false" json:"legal_hold"`
}
