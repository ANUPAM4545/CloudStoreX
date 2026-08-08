package model

import (
	"time"

	"github.com/google/uuid"
)

type PlanStatus string

const (
	PlanStatusActive    PlanStatus = "ACTIVE"
	PlanStatusDisabled  PlanStatus = "DISABLED"
	PlanStatusTriggered PlanStatus = "TRIGGERED"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

// RecoveryPlan defines an automated or manual disaster recovery procedure for a workspace.
type RecoveryPlan struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"workspace_id"`
	Name           string     `gorm:"type:varchar(255);not null" json:"name"`
	Description    string     `gorm:"type:text" json:"description"`
	SourceProvider string     `gorm:"type:varchar(100);not null" json:"source_provider"`
	TargetProvider string     `gorm:"type:varchar(100);not null" json:"target_provider"`
	RPOSeconds     int        `gorm:"not null;default:3600" json:"rpo_seconds"`
	RTOSeconds     int        `gorm:"not null;default:1800" json:"rto_seconds"`
	Status         PlanStatus `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
	LastTestedAt   *time.Time `json:"last_tested_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// RecoveryJob tracks the execution of a disaster recovery restore procedure (including Point-in-Time Recovery).
type RecoveryJob struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"plan_id"`
	WorkspaceID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"workspace_id"`
	SourceProvider  string     `gorm:"type:varchar(100);not null" json:"source_provider"`
	TargetProvider  string     `gorm:"type:varchar(100);not null" json:"target_provider"`
	PointInTime     *time.Time `json:"point_in_time"`
	Status          JobStatus  `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"status"`
	ObjectsRestored int64      `gorm:"not null;default:0" json:"objects_restored"`
	BytesRestored   int64      `gorm:"not null;default:0" json:"bytes_restored"`
	StartedAt       *time.Time `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	ErrorMessage    string     `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
