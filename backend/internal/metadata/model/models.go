package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bucket represents a logical container for objects in the metadata catalog.
type Bucket struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_bucket_name_workspace;index;not null"`
	ProviderID  string         `gorm:"type:varchar(100);index;not null"`
	Name        string         `gorm:"type:varchar(255);uniqueIndex:idx_bucket_name_workspace;not null"`
	Region      string         `gorm:"type:varchar(100)"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Object represents a stored object's metadata in the catalog.
// It maps the logical ObjectKey to the physical ProviderObjectKey.
type Object struct {
	ID                uuid.UUID      `gorm:"type:uuid;primaryKey"`
	BucketID          uuid.UUID      `gorm:"type:uuid;uniqueIndex:idx_bucket_object_key;index;not null"`
	ObjectKey         string         `gorm:"type:varchar(1024);uniqueIndex:idx_bucket_object_key;not null"`
	ProviderObjectKey string         `gorm:"type:varchar(1024);not null"`
	SizeBytes         int64          `gorm:"not null"`
	MimeType          string         `gorm:"type:varchar(255);index"`
	ETag              string         `gorm:"type:varchar(255);index"`
	ChecksumAlgorithm string         `gorm:"type:varchar(50)"`
	ProviderID        string         `gorm:"type:varchar(100);index;not null"`
	VersionID         *string        `gorm:"type:varchar(255)"`
	OwnerID           *uuid.UUID     `gorm:"type:uuid;index"`
	StorageClass      string         `gorm:"type:varchar(50);default:'STANDARD'"`
	Status            string         `gorm:"type:varchar(50);index;default:'ACTIVE'"` // ACTIVE, UPLOADING, ERROR
	IsDeleted         bool           `gorm:"default:false;index"`                     // Soft-delete flag
	CreatedAt         time.Time      `gorm:"autoCreateTime;index"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Relationships
	Bucket   Bucket           `gorm:"foreignKey:BucketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []ObjectTag      `gorm:"foreignKey:ObjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Metadata []ObjectMetadata `gorm:"foreignKey:ObjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// ObjectTag represents user-defined key-value tags for an object.
type ObjectTag struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ObjectID uuid.UUID `gorm:"type:uuid;index:idx_obj_tag_key,unique;not null"`
	Key      string    `gorm:"type:varchar(255);uniqueIndex:idx_obj_tag_key;index;not null"`
	Value    string    `gorm:"type:varchar(255);index"`
}

// ObjectMetadata represents system or user-defined custom metadata (e.g. headers).
type ObjectMetadata struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	ObjectID uuid.UUID `gorm:"type:uuid;index:idx_obj_meta_key,unique;not null"`
	Key      string    `gorm:"type:varchar(255);uniqueIndex:idx_obj_meta_key;index;not null"`
	Value    string    `gorm:"type:text"`
}
