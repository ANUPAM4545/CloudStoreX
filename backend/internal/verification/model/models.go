package model

import (
	"time"

	"github.com/google/uuid"
)

type CheckStatus string

const (
	StatusPending CheckStatus = "PENDING"
	StatusRunning CheckStatus = "RUNNING"
	StatusPassed  CheckStatus = "PASSED"
	StatusFailed  CheckStatus = "FAILED"
)

type DiscrepancyType string

const (
	DiscrepancyMissingInProvider DiscrepancyType = "MISSING_IN_PROVIDER"
	DiscrepancyMissingInMetadata DiscrepancyType = "MISSING_IN_METADATA"
	DiscrepancyChecksumMismatch  DiscrepancyType = "CHECKSUM_MISMATCH"
	DiscrepancySizeMismatch      DiscrepancyType = "SIZE_MISMATCH"
)

// ConsistencyJob tracks a scheduled or manual verification run over a workspace/bucket.
type ConsistencyJob struct {
	ID                  uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID         uuid.UUID   `gorm:"type:uuid;not null;index" json:"workspace_id"`
	BucketID            uuid.UUID   `gorm:"type:uuid;not null;index" json:"bucket_id"`
	ProviderID          string      `gorm:"type:varchar(100);not null" json:"provider_id"`
	Status              CheckStatus `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"status"`
	TotalObjectsChecked int64       `gorm:"not null;default:0" json:"total_objects_checked"`
	DiscrepanciesFound  int64       `gorm:"not null;default:0" json:"discrepancies_found"`
	StartedAt           *time.Time  `json:"started_at"`
	CompletedAt         *time.Time  `json:"completed_at"`
	ErrorMessage        string      `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt           time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

// DiscrepancyRecord details a specific consistency violation found during verification.
type DiscrepancyRecord struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	JobID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"job_id"`
	WorkspaceID uuid.UUID       `gorm:"type:uuid;not null;index" json:"workspace_id"`
	BucketID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"bucket_id"`
	ObjectID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"object_id"`
	ObjectKey   string          `gorm:"type:varchar(512);not null" json:"object_key"`
	Type        DiscrepancyType `gorm:"type:varchar(40);not null" json:"type"`
	Details     string          `gorm:"type:text" json:"details"`
	Resolved    bool            `gorm:"not null;default:false" json:"resolved"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
}
