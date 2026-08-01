package quota

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/quota/model"
)

var (
	ErrQuotaExceeded = errors.New("quota exceeded")
)

type Repository interface {
	GetWorkspaceQuota(ctx context.Context, workspaceID string) (*model.WorkspaceQuota, error)
	UpsertWorkspaceQuota(ctx context.Context, quota *model.WorkspaceQuota) error
	IncrementUsage(ctx context.Context, workspaceID string, bytes int64, objects int64) error
	DecrementUsage(ctx context.Context, workspaceID string, bytes int64, objects int64) error
}
