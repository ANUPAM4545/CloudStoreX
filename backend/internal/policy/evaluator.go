package policy

import (
	"context"

	"github.com/google/uuid"
)

// EvaluationContext encapsulates all metadata needed by the Policy Engine
// to make routing, compliance, replication, and cost-optimization decisions.
type EvaluationContext struct {
	WorkspaceID    uuid.UUID
	OrganizationID uuid.UUID
	Bucket         string
	ObjectKey      string
	Size           int64
	Region         string
	ComplianceTags []string
	CostOptimized  bool
	Metadata       map[string]string
}

// Evaluator defines the contract for evaluating storage policies.
type Evaluator interface {
	// EvaluateProvider decides which storage provider ID should handle a primary storage request.
	EvaluateProvider(ctx context.Context, evalCtx *EvaluationContext) (string, error)

	// EvaluateReplicationTargets returns a list of secondary provider IDs for replication.
	EvaluateReplicationTargets(ctx context.Context, evalCtx *EvaluationContext) ([]string, error)
}

// DefaultEvaluator is a lightweight placeholder implementation of Evaluator.
// It always returns a default provider identifier, preserving architecture for V1.
type DefaultEvaluator struct {
	defaultProviderID string
}

// NewDefaultEvaluator creates a new DefaultEvaluator.
func NewDefaultEvaluator(defaultProviderID string) *DefaultEvaluator {
	if defaultProviderID == "" {
		defaultProviderID = "default"
	}
	return &DefaultEvaluator{defaultProviderID: defaultProviderID}
}

func (e *DefaultEvaluator) EvaluateProvider(ctx context.Context, evalCtx *EvaluationContext) (string, error) {
	// For V1 MVP / placeholder, we always route to the default provider.
	return e.defaultProviderID, nil
}

func (e *DefaultEvaluator) EvaluateReplicationTargets(ctx context.Context, evalCtx *EvaluationContext) ([]string, error) {
	// For V1 MVP, no secondary replication targets are evaluated.
	return []string{}, nil
}
