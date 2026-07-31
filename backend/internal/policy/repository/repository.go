package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// Repository defines the contract for persisting policies and routing decisions.
type Repository interface {
	Create(ctx context.Context, p *model.Policy) error
	GetByID(ctx context.Context, id string) (*model.Policy, error)
	List(ctx context.Context, workspaceID string) ([]*model.Policy, error)
	ListEnabled(ctx context.Context, workspaceID string) ([]*model.Policy, error)
	Update(ctx context.Context, p *model.Policy) error
	Delete(ctx context.Context, id string) error

	// Audit trail
	CreateRoutingDecision(ctx context.Context, d *model.RoutingDecision) error
	ListRoutingDecisions(ctx context.Context, workspaceID string, limit, offset int) ([]*model.RoutingDecision, int64, error)
}
