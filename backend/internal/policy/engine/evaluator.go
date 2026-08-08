package engine

import (
	"context"
	"encoding/json"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// Evaluator evaluates an EvaluationContext against a set of policies to determine a target provider.
type Evaluator interface {
	Evaluate(ctx context.Context, evalCtx *model.EvaluationContext, policies []*model.Policy) (PolicyResult, *model.Policy, error)
}

// PolicyResult holds the outcome of evaluating a policy action
type PolicyResult struct {
	Matched             bool
	PrimaryProvider     string
	FallbackProviders   []string
	RejectIfUnsupported bool
	Explanation         string
}

type defaultEvaluator struct{}

// NewEvaluator creates a new Evaluation Engine.
func NewEvaluator() Evaluator {
	return &defaultEvaluator{}
}

func (e *defaultEvaluator) Evaluate(ctx context.Context, evalCtx *model.EvaluationContext, policies []*model.Policy) (PolicyResult, *model.Policy, error) {
	for _, p := range policies {
		if !p.Enabled {
			continue
		}

		var rootNode model.ConditionNode
		if len(p.Conditions) > 0 {
			if err := json.Unmarshal(p.Conditions, &rootNode); err != nil {
				// Malformed policy condition, skip
				continue
			}
		} else {
			// Empty condition implies match all
			rootNode = model.ConditionNode{} 
		}

		matched, err := evaluateNode(evalCtx, &rootNode)
		if err != nil {
			// Evaluation failed (e.g. type mismatch), skip policy
			continue
		}

		if matched {
			var action model.PolicyAction
			if len(p.Actions) > 0 {
				if err := json.Unmarshal(p.Actions, &action); err != nil {
					continue // Malformed action
				}
			}

			return PolicyResult{
				Matched:             true,
				PrimaryProvider:     action.PrimaryProvider,
				FallbackProviders:   action.FallbackProviders,
				RejectIfUnsupported: action.RejectIfUnsupported,
				Explanation:         "Policy matched based on dynamic rules.",
			}, p, nil
		}
	}

	// No policy matched
	return PolicyResult{Matched: false}, nil, nil
}
