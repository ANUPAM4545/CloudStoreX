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
	RuleTypeReplication  RuleType = "REPLICATION"
	RuleTypeDefault      RuleType = "DEFAULT"
)

// Policy defines a dynamic routing or orchestration policy.
type Policy struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Priority    int            `gorm:"not null;index"` // Lower number = higher priority
	Enabled     bool           `gorm:"default:true;not null;index"`
	RuleType    RuleType       `gorm:"type:varchar(50);not null"`
	Conditions  datatypes.JSON `gorm:"type:jsonb"` // Serialized ConditionNode
	Actions     datatypes.JSON `gorm:"type:jsonb"` // Serialized PolicyAction
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

// ConditionOperator defines the operation for a condition leaf
type ConditionOperator string

const (
	OpEq          ConditionOperator = "=="
	OpNeq         ConditionOperator = "!="
	OpContains    ConditionOperator = "contains"
	OpStartsWith  ConditionOperator = "starts_with"
	OpEndsWith    ConditionOperator = "ends_with"
	OpGt          ConditionOperator = ">"
	OpGte         ConditionOperator = ">="
	OpLt          ConditionOperator = "<"
	OpLte         ConditionOperator = "<="
	OpIn          ConditionOperator = "in"
	OpNotIn       ConditionOperator = "not_in"
	OpRegex       ConditionOperator = "regex"
	OpExists      ConditionOperator = "exists"
)

// ConditionNodeType defines if a node is an AND/OR group or a LEAF
type ConditionNodeType string

const (
	NodeAnd  ConditionNodeType = "AND"
	NodeOr   ConditionNodeType = "OR"
	NodeLeaf ConditionNodeType = "LEAF"
)

// ConditionLeaf holds the actual evaluation logic
type ConditionLeaf struct {
	Field    string            `json:"field"`
	Operator ConditionOperator `json:"operator"`
	Value    interface{}       `json:"value"`
}

// ConditionNode represents a node in the AST
type ConditionNode struct {
	Type       ConditionNodeType `json:"type"`
	Leaf       *ConditionLeaf    `json:"leaf,omitempty"`
	Conditions []ConditionNode   `json:"conditions,omitempty"`
}

// PolicyAction defines what happens when a policy matches
type PolicyAction struct {
	PrimaryProvider     string   `json:"primary_provider"`
	FallbackProviders   []string `json:"fallback_providers,omitempty"`
	RejectIfUnsupported bool     `json:"reject_if_unsupported"`
}

// RoutingDecision represents an auditable record of a policy routing decision.
type RoutingDecision struct {
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkspaceID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	ObjectKey          string         `gorm:"type:varchar(1024);index"`
	ProviderID         string         `gorm:"type:varchar(255);index"`
	PolicyID           *uuid.UUID     `gorm:"type:uuid;index"`
	MatchedRule        string         `gorm:"type:varchar(255)"`
	IsFallback         bool           `gorm:"default:false"`
	CandidateProviders datatypes.JSON `gorm:"type:jsonb"` // Array of EvaluationRecord logs
	DecisionReason     string         `gorm:"type:text"`
	Operation          string         `gorm:"type:varchar(50)"`
	LatencyMs          int64          `gorm:"not null"`
	Timestamp          time.Time      `gorm:"not null;index"`
}

// EvaluationRecord captures why a provider was accepted or rejected during fallback
type EvaluationRecord struct {
	ProviderID string `json:"provider_id"`
	Status     string `json:"status"` // "SELECTED" or "REJECTED"
	Reason     string `json:"reason"`
}
