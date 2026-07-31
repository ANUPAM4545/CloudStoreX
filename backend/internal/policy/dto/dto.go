package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PolicyDTO struct {
	ID          string          `json:"id"`
	WorkspaceID string          `json:"workspace_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Priority    int             `json:"priority"`
	Enabled     bool            `json:"enabled"`
	RuleType    string          `json:"rule_type"`
	Conditions  json.RawMessage `json:"conditions"`
	Actions     json.RawMessage `json:"actions"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CreatePolicyRequest struct {
	WorkspaceID uuid.UUID       `json:"workspace_id"`
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Priority    int             `json:"priority"`
	RuleType    string          `json:"rule_type" binding:"required"`
	Conditions  json.RawMessage `json:"conditions"`
	Actions     json.RawMessage `json:"actions" binding:"required"`
}

type UpdatePolicyRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Priority    int             `json:"priority"`
	RuleType    string          `json:"rule_type"`
	Conditions  json.RawMessage `json:"conditions"`
	Actions     json.RawMessage `json:"actions"`
}

type RoutingDecisionDTO struct {
	ID                 string          `json:"id"`
	WorkspaceID        string          `json:"workspace_id"`
	ObjectKey          string          `json:"object_key"`
	ProviderID         string          `json:"provider_id"`
	PolicyID           *string         `json:"policy_id,omitempty"`
	MatchedRule        string          `json:"matched_rule"`
	IsFallback         bool            `json:"is_fallback"`
	CandidateProviders json.RawMessage `json:"candidate_providers"`
	DecisionReason     string          `json:"decision_reason"`
	Operation          string          `json:"operation"`
	LatencyMs          int64           `json:"latency_ms"`
	Timestamp          time.Time       `json:"timestamp"`
}
