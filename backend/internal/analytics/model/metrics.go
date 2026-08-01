package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StorageSnapshot captures historical storage utilization for a workspace.
type StorageSnapshot struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	WorkspaceID string         `gorm:"type:varchar(255);index;not null" json:"workspace_id"`
	TotalBytes  int64          `gorm:"not null;default:0" json:"total_bytes"`
	TotalObjects int64         `gorm:"not null;default:0" json:"total_objects"`
	BucketCount int64          `gorm:"not null;default:0" json:"bucket_count"`
	Date        time.Time      `gorm:"type:date;index;not null" json:"date"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
