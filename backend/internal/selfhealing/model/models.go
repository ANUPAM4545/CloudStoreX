package model

import (
	"time"

	"github.com/google/uuid"
)

type RepairStatus string

const (
	StatusPending   RepairStatus = "PENDING"
	StatusRunning   RepairStatus = "RUNNING"
	StatusCompleted RepairStatus = "COMPLETED"
	StatusFailed    RepairStatus = "FAILED"
)

type RepairSource string

const (
	SourceReplica RepairSource = "HEALTHY_REPLICA"
	SourceBackup  RepairSource = "BACKUP_VAULT"
)

// RepairJob tracks an automatic or manual self-healing repair execution for an object discrepancy.
type RepairJob struct {
	ID             uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"workspace_id"`
	BucketID       uuid.UUID    `gorm:"type:uuid;not null;index" json:"bucket_id"`
	ObjectID       uuid.UUID    `gorm:"type:uuid;not null;index" json:"object_id"`
	ObjectKey      string       `gorm:"type:varchar(512);not null" json:"object_key"`
	TargetProvider string       `gorm:"type:varchar(100);not null" json:"target_provider"`
	SourceProvider string       `gorm:"type:varchar(100);not null" json:"source_provider"`
	SourceType     RepairSource `gorm:"type:varchar(30);not null;default:'HEALTHY_REPLICA'" json:"source_type"`
	Status         RepairStatus `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"status"`
	BytesRepaired  int64        `gorm:"not null;default:0" json:"bytes_repaired"`
	StartedAt      *time.Time   `json:"started_at"`
	CompletedAt    *time.Time   `json:"completed_at"`
	ErrorMessage   string       `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}
