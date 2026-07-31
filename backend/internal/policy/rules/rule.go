package rules

import (
	"encoding/json"
	"fmt"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// RuleResult represents the outcome of a rule evaluation.
type RuleResult struct {
	Matched     bool
	ProviderID  string
	Explanation string
}

// Rule defines the interface that all policy rules must implement.
type Rule interface {
	// RuleType returns the type identifier for this rule.
	RuleType() model.RuleType
	
	// Evaluate determines if the policy conditions match the evaluation context
	// and returns the outcome (e.g., Target Provider).
	Evaluate(ctx *model.EvaluationContext, conditionsJSON, actionsJSON []byte) (RuleResult, error)
}

// Actions represents the common structure for rule actions.
type Actions struct {
	ProviderID string `json:"provider_id"`
}

// parseActions is a helper to parse the common action structure.
func parseActions(actionsJSON []byte) (*Actions, error) {
	var actions Actions
	if len(actionsJSON) == 0 {
		return &actions, nil
	}
	if err := json.Unmarshal(actionsJSON, &actions); err != nil {
		return nil, fmt.Errorf("failed to parse actions: %w", err)
	}
	return &actions, nil
}
