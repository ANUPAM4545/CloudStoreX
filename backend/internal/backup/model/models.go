package model

import (
	"time"

	"github.com/google/uuid"
)

type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "PENDING"
	BackupStatusRunning   BackupStatus = "RUNNING"
	BackupStatusCompleted BackupStatus = "COMPLETED"
	BackupStatusFailed    BackupStatus = "FAILED"
)

type BackupType string

const (
	BackupTypeFull        BackupType = "FULL"
	BackupTypeIncremental BackupType = "INCREMENTAL"
)

// BackupJob tracks a scheduled or manual backup execution.
type BackupJob struct {
	ID              uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID     uuid.UUID    `gorm:"type:uuid;not null;index" json:"workspace_id"`
	BucketID        uuid.UUID    `gorm:"type:uuid;not null;index" json:"bucket_id"`
	SourceProvider  string       `gorm:"type:varchar(100);not null" json:"source_provider"`
	TargetProvider  string       `gorm:"type:varchar(100);not null" json:"target_provider"`
	BackupType      BackupType   `gorm:"type:varchar(20);not null;default:'FULL'" json:"backup_type"`
	Status          BackupStatus `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"status"`
	ObjectsBackedUp int64        `gorm:"not null;default:0" json:"objects_backed_up"`
	BytesBackedUp   int64        `gorm:"not null;default:0" json:"bytes_backed_up"`
	SnapshotID      string       `gorm:"type:varchar(128)" json:"snapshot_id"`
	StartedAt       *time.Time   `json:"started_at"`
	CompletedAt     *time.Time   `json:"completed_at"`
	ErrorMessage    string       `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

// BackupRecord represents an immutable record of an object or snapshot stored in backup vault storage.
type BackupRecord struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BackupJobID   uuid.UUID `gorm:"type:uuid;not null;index" json:"backup_job_id"`
	WorkspaceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"workspace_id"`
	BucketID      uuid.UUID `gorm:"type:uuid;not null;index" json:"bucket_id"`
	SnapshotID    string    `gorm:"type:varchar(128);not null;index" json:"snapshot_id"`
	VaultLocation string    `gorm:"type:varchar(512);not null" json:"vault_location"`
	TotalBytes    int64     `gorm:"not null;default:0" json:"total_bytes"`
	Checksum      string    `gorm:"type:varchar(64)" json:"checksum"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
