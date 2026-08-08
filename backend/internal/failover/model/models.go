package model

import (
	"time"

	"github.com/google/uuid"
)

type FailoverStatus string

const (
	StatusActive    FailoverStatus = "ACTIVE"
	StatusCompleted FailoverStatus = "COMPLETED"
	StatusFailed    FailoverStatus = "FAILED"
	StatusFailback  FailoverStatus = "FAILBACK_COMPLETED"
)

// FailoverEvent records an automatic or manual failover transition between storage providers.
type FailoverEvent struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"workspace_id"`
	OriginalProvider   string         `gorm:"type:varchar(100);not null" json:"original_provider"`
	NewPrimaryProvider string         `gorm:"type:varchar(100);not null" json:"new_primary_provider"`
	Reason             string         `gorm:"type:text;not null" json:"reason"`
	Status             FailoverStatus `gorm:"type:varchar(30);not null;default:'ACTIVE'" json:"status"`
	TriggeredAt        time.Time      `gorm:"autoCreateTime" json:"triggered_at"`
	ResolvedAt         *time.Time     `json:"resolved_at,omitempty"`
}
