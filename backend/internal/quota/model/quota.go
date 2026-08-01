package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WorkspaceQuota represents storage capacity limits and usage for a workspace.
type WorkspaceQuota struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	WorkspaceID string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	MaxBytes    int64          `gorm:"not null;default:0"` // 0 means unlimited
	BytesUsed   int64          `gorm:"not null;default:0"`
	MaxObjects  int64          `gorm:"not null;default:0"` // 0 means unlimited
	ObjectCount int64          `gorm:"not null;default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
