package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/provider/model"
)

// Repository defines the data access contract for Provider entities.
type Repository interface {
	Create(ctx context.Context, provider *model.Provider) error
	GetByID(ctx context.Context, id string) (*model.Provider, error)
	GetByName(ctx context.Context, workspaceID, name string) (*model.Provider, error)
	List(ctx context.Context, workspaceID string) ([]*model.Provider, error)
	GetDefault(ctx context.Context, workspaceID string) (*model.Provider, error)
	Update(ctx context.Context, provider *model.Provider) error
	UpdateHealth(ctx context.Context, id string, status model.ProviderStatus, health model.ProviderHealth, latencyMs int64, lastError string) error
	SetDefault(ctx context.Context, workspaceID, id string) error
	Enable(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
