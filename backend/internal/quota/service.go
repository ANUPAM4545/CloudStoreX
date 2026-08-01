package quota

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cloudstorex/backend/internal/quota/model"
)

type Service interface {
	CheckQuota(ctx context.Context, workspaceID string, sizeBytes int64) error
	IncrementUsage(ctx context.Context, workspaceID string, sizeBytes int64, objects int64) error
	DecrementUsage(ctx context.Context, workspaceID string, sizeBytes int64, objects int64) error
	SetWorkspaceQuota(ctx context.Context, quota *model.WorkspaceQuota) error
	GetWorkspaceQuota(ctx context.Context, workspaceID string) (*model.WorkspaceQuota, error)
}

type defaultService struct {
	repo Repository
	log  *slog.Logger
}

func NewService(repo Repository, log *slog.Logger) Service {
	if log == nil {
		log = slog.Default()
	}
	return &defaultService{
		repo: repo,
		log:  log,
	}
}

func (s *defaultService) CheckQuota(ctx context.Context, workspaceID string, sizeBytes int64) error {
	q, err := s.repo.GetWorkspaceQuota(ctx, workspaceID)
	if err != nil {
		// If no quota is set, assume unlimited or handle missing quota based on policy.
		// For now, we'll allow it if missing, but we could return an error.
		return nil
	}

	if q.MaxBytes > 0 && q.BytesUsed+sizeBytes > q.MaxBytes {
		s.log.WarnContext(ctx, "quota exceeded (bytes)",
			slog.String("workspace_id", workspaceID),
			slog.Int64("requested", sizeBytes),
			slog.Int64("used", q.BytesUsed),
			slog.Int64("max", q.MaxBytes),
		)
		return ErrQuotaExceeded
	}

	if q.MaxObjects > 0 && q.ObjectCount+1 > q.MaxObjects {
		s.log.WarnContext(ctx, "quota exceeded (objects)",
			slog.String("workspace_id", workspaceID),
			slog.Int64("used", q.ObjectCount),
			slog.Int64("max", q.MaxObjects),
		)
		return ErrQuotaExceeded
	}

	return nil
}

func (s *defaultService) IncrementUsage(ctx context.Context, workspaceID string, sizeBytes int64, objects int64) error {
	return s.repo.IncrementUsage(ctx, workspaceID, sizeBytes, objects)
}

func (s *defaultService) DecrementUsage(ctx context.Context, workspaceID string, sizeBytes int64, objects int64) error {
	return s.repo.DecrementUsage(ctx, workspaceID, sizeBytes, objects)
}

func (s *defaultService) SetWorkspaceQuota(ctx context.Context, quota *model.WorkspaceQuota) error {
	if quota.WorkspaceID == "" {
		return fmt.Errorf("workspace ID is required")
	}
	return s.repo.UpsertWorkspaceQuota(ctx, quota)
}

func (s *defaultService) GetWorkspaceQuota(ctx context.Context, workspaceID string) (*model.WorkspaceQuota, error) {
	return s.repo.GetWorkspaceQuota(ctx, workspaceID)
}
