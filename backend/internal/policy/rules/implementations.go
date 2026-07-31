package rules

import (
	"encoding/json"
	"fmt"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// ObjectSizeRule routes objects based on their byte size.
type ObjectSizeRule struct{}

type sizeConditions struct {
	MinSize int64 `json:"min_size"`
	MaxSize int64 `json:"max_size"`
}

func (r *ObjectSizeRule) RuleType() model.RuleType {
	return model.RuleTypeObjectSize
}

func (r *ObjectSizeRule) Evaluate(ctx *model.EvaluationContext, conditionsJSON, actionsJSON []byte) (RuleResult, error) {
	var conds sizeConditions
	if len(conditionsJSON) > 0 {
		if err := json.Unmarshal(conditionsJSON, &conds); err != nil {
			return RuleResult{}, fmt.Errorf("failed to parse size conditions: %w", err)
		}
	}

	actions, err := parseActions(actionsJSON)
	if err != nil {
		return RuleResult{}, err
	}

	matches := true
	if conds.MinSize > 0 && ctx.Size < conds.MinSize {
		matches = false
	}
	if conds.MaxSize > 0 && ctx.Size > conds.MaxSize {
		matches = false
	}

	if matches {
		return RuleResult{
			Matched:     true,
			ProviderID:  actions.ProviderID,
			Explanation: fmt.Sprintf("Object size %d matched constraints (min: %d, max: %d)", ctx.Size, conds.MinSize, conds.MaxSize),
		}, nil
	}

	return RuleResult{Matched: false}, nil
}

// RegionRule routes objects based on region preference.
type RegionRule struct{}

type regionConditions struct {
	Region string `json:"region"`
}

func (r *RegionRule) RuleType() model.RuleType {
	return model.RuleTypeRegion
}

func (r *RegionRule) Evaluate(ctx *model.EvaluationContext, conditionsJSON, actionsJSON []byte) (RuleResult, error) {
	var conds regionConditions
	if len(conditionsJSON) > 0 {
		if err := json.Unmarshal(conditionsJSON, &conds); err != nil {
			return RuleResult{}, fmt.Errorf("failed to parse region conditions: %w", err)
		}
	}

	actions, err := parseActions(actionsJSON)
	if err != nil {
		return RuleResult{}, err
	}

	if ctx.Region != "" && conds.Region != "" && ctx.Region == conds.Region {
		return RuleResult{
			Matched:     true,
			ProviderID:  actions.ProviderID,
			Explanation: fmt.Sprintf("Context region %s matched condition", ctx.Region),
		}, nil
	}

	return RuleResult{Matched: false}, nil
}

// DefaultRule serves as a simple explicit default provider assignment.
type DefaultRule struct{}

func (r *DefaultRule) RuleType() model.RuleType {
	return model.RuleTypeDefault
}

func (r *DefaultRule) Evaluate(ctx *model.EvaluationContext, conditionsJSON, actionsJSON []byte) (RuleResult, error) {
	actions, err := parseActions(actionsJSON)
	if err != nil {
		return RuleResult{}, err
	}
	return RuleResult{
		Matched:     true,
		ProviderID:  actions.ProviderID,
		Explanation: "Default rule applied",
	}, nil
}
