package engine

import (
	"context"
	"time"

	"github.com/cloudstorex/backend/internal/policy/events"
	"github.com/cloudstorex/backend/internal/policy/model"
)

// PolicyRepository defines the subset of repository methods needed by the engine.
type PolicyRepository interface {
	ListEnabled(ctx context.Context, workspaceID string) ([]*model.Policy, error)
}

// DefaultProviderResolver allows fetching the fallback provider.
type DefaultProviderResolver interface {
	GetDefaultProviderID(ctx context.Context, workspaceID string) (string, error)
}

// PolicyEngine orchestrates dynamic routing decisions.
type PolicyEngine interface {
	Resolve(ctx context.Context, evalCtx *model.EvaluationContext, op string) (string, *model.RoutingDecision, error)
}

type defaultPolicyEngine struct {
	repo         PolicyRepository
	evaluator    Evaluator
	resolver     DefaultProviderResolver
	eventPub     events.Publisher
}

// NewPolicyEngine creates a new Policy Engine instance.
func NewPolicyEngine(repo PolicyRepository, evaluator Evaluator, resolver DefaultProviderResolver, eventPub events.Publisher) PolicyEngine {
	return &defaultPolicyEngine{
		repo:      repo,
		evaluator: evaluator,
		resolver:  resolver,
		eventPub:  eventPub,
	}
}

func (e *defaultPolicyEngine) Resolve(ctx context.Context, evalCtx *model.EvaluationContext, op string) (string, *model.RoutingDecision, error) {
	start := time.Now()
	
	decision := &model.RoutingDecision{
		WorkspaceID: evalCtx.WorkspaceID,
		ObjectKey:   evalCtx.ObjectKey,
		Operation:   op,
		Timestamp:   start,
	}
	
	var finalProviderID string

	// 1. Load active policies
	policies, err := e.repo.ListEnabled(ctx, evalCtx.WorkspaceID.String())
	if err != nil {
		// Log error, but we can fallback
		decision.DecisionReason = "Failed to load policies, using fallback."
	} else {
		// 2. Evaluate rules
		result, matchedPolicy, _ := e.evaluator.Evaluate(ctx, evalCtx, policies)
		
		if result.Matched {
			finalProviderID = result.ProviderID
			decision.PolicyID = &matchedPolicy.ID
			decision.MatchedRule = string(matchedPolicy.RuleType)
			decision.DecisionReason = result.Explanation
			decision.IsFallback = false
		}
	}

	// 3. Fallback
	if finalProviderID == "" {
		decision.IsFallback = true
		fallbackID, fErr := e.resolver.GetDefaultProviderID(ctx, evalCtx.WorkspaceID.String())
		if fErr != nil || fallbackID == "" {
			// Absolute last resort
			fallbackID = "minio"
			decision.DecisionReason = "No policy matched and default provider resolution failed. Absolute fallback applied."
		} else {
			if decision.DecisionReason == "" {
				decision.DecisionReason = "No policy matched. Fallback to default provider applied."
			}
		}
		finalProviderID = fallbackID
	}

	decision.ProviderID = finalProviderID
	decision.LatencyMs = time.Since(start).Milliseconds()

	// 4. Async Audit Log Emission
	_ = e.eventPub.PublishRoutingDecision(context.Background(), decision)

	return finalProviderID, decision, nil
}
