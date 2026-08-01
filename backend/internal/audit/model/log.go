package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog represents an immutable record of a significant system or user action.
type AuditLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	WorkspaceID string         `gorm:"type:varchar(255);index;not null" json:"workspace_id"`
	UserID      string         `gorm:"type:varchar(255);index" json:"user_id,omitempty"`
	Action      string         `gorm:"type:varchar(100);index;not null" json:"action"`
	Resource    string         `gorm:"type:varchar(255);index" json:"resource,omitempty"`
	Details     datatypes.JSON `gorm:"type:jsonb" json:"details,omitempty"`
	Timestamp   time.Time      `gorm:"autoCreateTime;index" json:"timestamp"`
}
