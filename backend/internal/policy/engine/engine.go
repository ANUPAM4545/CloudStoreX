package engine

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/cloudstorex/backend/internal/observability/metrics"
	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/cloudstorex/backend/internal/policy/events"
	"github.com/cloudstorex/backend/internal/policy/model"
	"go.opentelemetry.io/otel/attribute"
)

// PolicyRepository defines the subset of repository methods needed by the engine.
type PolicyRepository interface {
	ListEnabled(ctx context.Context, workspaceID string) ([]*model.Policy, error)
}

// DefaultProviderResolver allows fetching the fallback provider.
type DefaultProviderResolver interface {
	GetDefaultProviderID(ctx context.Context, workspaceID string) (string, error)
}

// ProviderValidator validates if a provider is ready to handle the request
type ProviderValidator interface {
	// Validate returns healthy, capable, error.
	// op is the operation name e.g. "UploadObject", "CopyObject"
	Validate(ctx context.Context, providerID string, op string) (healthy bool, capable bool, err error)
}

// PolicyEngine orchestrates dynamic routing decisions.
type PolicyEngine interface {
	Resolve(ctx context.Context, evalCtx *model.EvaluationContext, op string) (string, *model.RoutingDecision, error)
}

type defaultPolicyEngine struct {
	repo         PolicyRepository
	evaluator    Evaluator
	resolver     DefaultProviderResolver
	validator    ProviderValidator
	eventPub     events.Publisher
}

// NewPolicyEngine creates a new Policy Engine instance.
func NewPolicyEngine(repo PolicyRepository, evaluator Evaluator, resolver DefaultProviderResolver, validator ProviderValidator, eventPub events.Publisher) PolicyEngine {
	return &defaultPolicyEngine{
		repo:      repo,
		evaluator: evaluator,
		resolver:  resolver,
		validator: validator,
		eventPub:  eventPub,
	}
}

func (e *defaultPolicyEngine) Resolve(ctx context.Context, evalCtx *model.EvaluationContext, op string) (string, *model.RoutingDecision, error) {
	ctx, span := tracing.StartChildSpan(ctx, "PolicyEngine.Resolve")
	defer span.End()

	start := time.Now()

	decision := &model.RoutingDecision{
		WorkspaceID: evalCtx.WorkspaceID,
		ObjectKey:   evalCtx.ObjectKey,
		Operation:   op,
		Timestamp:   start,
	}
	
	// Increment total evaluations
	metrics.PolicyEvaluationsTotal.WithLabelValues(evalCtx.WorkspaceID.String(), op, "started").Inc()

	var finalProviderID string
	var candidates []model.EvaluationRecord

	// 1. Load active policies
	policies, err := e.repo.ListEnabled(ctx, evalCtx.WorkspaceID.String())
	if err != nil {
		candidates = append(candidates, model.EvaluationRecord{
			ProviderID: "none",
			Status:     "REJECTED",
			Reason:     "Failed to load policies",
		})
		metrics.PolicyEvaluationsTotal.WithLabelValues(evalCtx.WorkspaceID.String(), op, "error").Inc()
	} else {
		// 2. Sort policies deterministically: Priority ASC, then ID ASC
		sort.SliceStable(policies, func(i, j int) bool {
			if policies[i].Priority != policies[j].Priority {
				return policies[i].Priority < policies[j].Priority
			}
			return policies[i].ID.String() < policies[j].ID.String()
		})

		// 3. Evaluate rules
		result, matchedPolicy, _ := e.evaluator.Evaluate(ctx, evalCtx, policies)
		
		if result.Matched {
			decision.PolicyID = &matchedPolicy.ID
			decision.MatchedRule = string(matchedPolicy.RuleType)
			
			// Try primary provider
			candidates = append(candidates, e.evaluateCandidate(ctx, result.PrimaryProvider, op, true))
			
			if candidates[len(candidates)-1].Status == "SELECTED" {
				finalProviderID = result.PrimaryProvider
				decision.DecisionReason = "Primary provider validated successfully."
				decision.IsFallback = false
			} else {
				// Try fallbacks
				for _, fallbackProvider := range result.FallbackProviders {
					candidates = append(candidates, e.evaluateCandidate(ctx, fallbackProvider, op, false))
					if candidates[len(candidates)-1].Status == "SELECTED" {
						finalProviderID = fallbackProvider
						decision.DecisionReason = "Fell back to alternative provider: " + candidates[len(candidates)-1].Reason
						decision.IsFallback = true
						metrics.PolicyFallbackRoutes.WithLabelValues(evalCtx.WorkspaceID.String(), op).Inc()
						break
					}
				}
				
				if finalProviderID == "" && result.RejectIfUnsupported {
					decision.DecisionReason = "Policy mandates rejection if primary and fallbacks are unsupported."
					metrics.PolicyEvaluationsTotal.WithLabelValues(evalCtx.WorkspaceID.String(), op, "rejected").Inc()
				}
			}
		}
	}

	// 4. Default Fallback
	if finalProviderID == "" && decision.DecisionReason == "" {
		decision.IsFallback = true
		fallbackID, fErr := e.resolver.GetDefaultProviderID(ctx, evalCtx.WorkspaceID.String())
		
		if fErr == nil && fallbackID != "" {
			candidates = append(candidates, e.evaluateCandidate(ctx, fallbackID, op, false))
			if candidates[len(candidates)-1].Status == "SELECTED" {
				finalProviderID = fallbackID
				decision.DecisionReason = "No policy matched or all candidates failed. Fallback to default provider applied."
				metrics.PolicyFallbackRoutes.WithLabelValues(evalCtx.WorkspaceID.String(), op).Inc()
			}
		}
		
		// Absolute last resort
		if finalProviderID == "" {
			finalProviderID = "minio"
			decision.DecisionReason = "Absolute fallback applied."
		}
	}

	decision.ProviderID = finalProviderID
	decision.LatencyMs = time.Since(start).Milliseconds()
	
	candidatesJSON, _ := json.Marshal(candidates)
	decision.CandidateProviders = candidatesJSON

	span.SetAttributes(
		attribute.String("policy.workspace_id", evalCtx.WorkspaceID.String()),
		attribute.String("policy.operation", op),
		attribute.String("policy.provider_id", finalProviderID),
		attribute.Bool("policy.is_fallback", decision.IsFallback),
	)
	
	if decision.MatchedRule != "" {
		metrics.PolicyRoutingDecisions.WithLabelValues(evalCtx.WorkspaceID.String(), finalProviderID, op, decision.MatchedRule).Inc()
	}

	// 5. Async Audit Log Emission
	_ = e.eventPub.PublishRoutingDecision(context.Background(), decision)

	metrics.PolicyEvaluationsTotal.WithLabelValues(evalCtx.WorkspaceID.String(), op, "success").Inc()

	return finalProviderID, decision, nil
}

func (e *defaultPolicyEngine) evaluateCandidate(ctx context.Context, providerID, op string, isPrimary bool) model.EvaluationRecord {
	rec := model.EvaluationRecord{
		ProviderID: providerID,
		Status:     "REJECTED",
	}
	
	healthy, capable, err := e.validator.Validate(ctx, providerID, op)
	if err != nil {
		rec.Reason = "validation error: " + err.Error()
		return rec
	}
	
	if !healthy {
		rec.Reason = "provider unhealthy"
		return rec
	}
	
	if !capable {
		rec.Reason = "missing required capability"
		return rec
	}
	
	rec.Status = "SELECTED"
	if isPrimary {
		rec.Reason = "healthy + capability compatible"
	} else {
		rec.Reason = "fallback selected (healthy + capability compatible)"
	}
	
	return rec
}
