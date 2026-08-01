package audit

import (
	"context"
	"log/slog"

	"github.com/cloudstorex/backend/internal/audit/model"
)

type Service interface {
	Record(ctx context.Context, log *model.AuditLog) error
	ListLogs(ctx context.Context, workspaceID string, limit, offset int) ([]model.AuditLog, int64, error)
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

func (s *defaultService) Record(ctx context.Context, logEntry *model.AuditLog) error {
	if err := s.repo.Create(ctx, logEntry); err != nil {
		s.log.ErrorContext(ctx, "failed to record audit log",
			slog.String("action", logEntry.Action),
			slog.String("workspace_id", logEntry.WorkspaceID),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (s *defaultService) ListLogs(ctx context.Context, workspaceID string, limit, offset int) ([]model.AuditLog, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListByWorkspace(ctx, workspaceID, limit, offset)
}
