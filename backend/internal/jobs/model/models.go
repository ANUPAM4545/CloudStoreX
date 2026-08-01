package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type JobState string

const (
	JobStatePending    JobState = "PENDING"
	JobStateProcessing JobState = "PROCESSING"
	JobStateCompleted  JobState = "COMPLETED"
	JobStateFailed     JobState = "FAILED"
	JobStateRetry      JobState = "RETRY"
)

// Job represents a background task that needs to be executed asynchronously.
type Job struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type        string         `gorm:"type:varchar(100);not null;index" json:"type"`
	Payload     datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	State       JobState       `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"state"`
	MaxRetries  int            `gorm:"not null;default:3" json:"max_retries"`
	RetryCount  int            `gorm:"not null;default:0" json:"retry_count"`
	Error       string         `gorm:"type:text" json:"error"`
	ScheduledAt *time.Time     `gorm:"index" json:"scheduled_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// JobLog records the execution history of jobs for debugging and audit.
type JobLog struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	JobID     uuid.UUID `gorm:"type:uuid;not null;index" json:"job_id"`
	State     JobState  `gorm:"type:varchar(20);not null" json:"state"`
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
