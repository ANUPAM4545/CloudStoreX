package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RuleType string

const (
	RuleTypeRegion       RuleType = "REGION"
	RuleTypeHealth       RuleType = "HEALTH"
	RuleTypeCapability   RuleType = "CAPABILITY"
	RuleTypeBucket       RuleType = "BUCKET"
	RuleTypeObjectSize   RuleType = "OBJECT_SIZE"
	RuleTypeMimeType     RuleType = "MIME_TYPE"
	RuleTypeDefault      RuleType = "DEFAULT"
)

// Policy defines a storage routing or orchestration policy.
type Policy struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Priority    int            `gorm:"not null"` // Lower number = higher priority
	Enabled     bool           `gorm:"default:true;not null"`
	RuleType    RuleType       `gorm:"type:varchar(50);not null"`
	Conditions  datatypes.JSON `gorm:"type:jsonb"` // JSON representation of the rule condition
	Actions     datatypes.JSON `gorm:"type:jsonb"` // JSON representation of the action (e.g., {"provider_id": "..."})
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

// RoutingDecision represents an auditable record of a policy routing decision.
type RoutingDecision struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkspaceID        uuid.UUID `gorm:"type:uuid;not null;index"`
	ObjectKey          string    `gorm:"type:varchar(1024);index"`
	ProviderID         string    `gorm:"type:varchar(255);index"`
	PolicyID           *uuid.UUID `gorm:"type:uuid;index"` // Nullable if fallback was used
	MatchedRule        string    `gorm:"type:varchar(255)"`
	IsFallback         bool      `gorm:"default:false"`
	CandidateProviders datatypes.JSON `gorm:"type:jsonb"` // Array of provider IDs considered
	DecisionReason     string    `gorm:"type:text"`
	Operation          string    `gorm:"type:varchar(50)"`
	LatencyMs          int64     `gorm:"not null"`
	Timestamp          time.Time `gorm:"not null;index"`
}
