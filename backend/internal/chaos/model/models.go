package model

import (
	"time"

	"github.com/google/uuid"
)

type ExperimentType string

const (
	TypeProviderOutage   ExperimentType = "PROVIDER_OUTAGE"
	TypeLatencyInjection ExperimentType = "LATENCY_INJECTION"
	TypeNetworkPartition ExperimentType = "NETWORK_PARTITION"
	TypeDataCorruption   ExperimentType = "DATA_CORRUPTION"
	TypeMetadataLag      ExperimentType = "METADATA_LAG"
	TypeWorkerCrash      ExperimentType = "WORKER_CRASH"
	TypeSplitBrain       ExperimentType = "SPLIT_BRAIN"
	TypeQuotaExhaustion  ExperimentType = "QUOTA_EXHAUSTION"
	TypeWriteCollision   ExperimentType = "WRITE_COLLISION"
)

type ExperimentStatus string

const (
	StatusScheduled ExperimentStatus = "SCHEDULED"
	StatusRunning   ExperimentStatus = "RUNNING"
	StatusCompleted ExperimentStatus = "COMPLETED"
	StatusFailed    ExperimentStatus = "FAILED"
	StatusAborted   ExperimentStatus = "ABORTED"
)

// ChaosExperiment defines a controlled reliability or fault-injection test scenario.
type ChaosExperiment struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"workspace_id"`
	Name            string           `gorm:"type:varchar(255);not null" json:"name"`
	ExperimentType  ExperimentType   `gorm:"type:varchar(50);not null" json:"experiment_type"`
	TargetProvider  string           `gorm:"type:varchar(100)" json:"target_provider,omitempty"`
	DurationSeconds int              `gorm:"not null;default:60" json:"duration_seconds"`
	Status          ExperimentStatus `gorm:"type:varchar(20);not null;default:'SCHEDULED';index" json:"status"`
	MeasuredRPO     int              `json:"measured_rpo,omitempty"` // in seconds
	MeasuredRTO     int              `json:"measured_rto,omitempty"` // in seconds
	TargetRPO       int              `json:"target_rpo,omitempty"`   // expected RPO objective
	TargetRTO       int              `json:"target_rto,omitempty"`   // expected RTO objective
	Passed          bool             `gorm:"not null;default:false" json:"passed"`
	StartedAt       *time.Time       `json:"started_at"`
	CompletedAt     *time.Time       `json:"completed_at"`
	ErrorMessage    string           `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}
