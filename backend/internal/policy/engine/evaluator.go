package engine

import (
	"context"
	"sort"

	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/cloudstorex/backend/internal/policy/rules"
)

// Evaluator evaluates an EvaluationContext against a set of policies to determine a target provider.
type Evaluator interface {
	Evaluate(ctx context.Context, evalCtx *model.EvaluationContext, policies []*model.Policy) (rules.RuleResult, *model.Policy, error)
}

type defaultEvaluator struct {
	registry rules.Registry
}

// NewEvaluator creates a new Evaluation Engine.
func NewEvaluator(registry rules.Registry) Evaluator {
	return &defaultEvaluator{
		registry: registry,
	}
}

func (e *defaultEvaluator) Evaluate(ctx context.Context, evalCtx *model.EvaluationContext, policies []*model.Policy) (rules.RuleResult, *model.Policy, error) {
	// Filter enabled policies
	var activePolicies []*model.Policy
	for _, p := range policies {
		if p.Enabled {
			activePolicies = append(activePolicies, p)
		}
	}

	// Sort by priority ASC (lower number is higher priority)
	sort.Slice(activePolicies, func(i, j int) bool {
		return activePolicies[i].Priority < activePolicies[j].Priority
	})

	for _, p := range activePolicies {
		rule, err := e.registry.Get(p.RuleType)
		if err != nil {
			// Skip unsupported rules to prevent pipeline blocking, log it in real scenario
			continue
		}

		result, err := rule.Evaluate(evalCtx, p.Conditions, p.Actions)
		if err != nil {
			// Rule evaluation failed, skip
			continue
		}

		if result.Matched {
			return result, p, nil
		}
	}

	// No policy matched
	return rules.RuleResult{Matched: false}, nil, nil
}
