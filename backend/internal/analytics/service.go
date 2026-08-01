package analytics

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/analytics/model"
)

type Service interface {
	RecordSnapshot(ctx context.Context, snapshot *model.StorageSnapshot) error
	GetHistory(ctx context.Context, workspaceID string, from, to time.Time) ([]model.StorageSnapshot, error)
	GetLatestMetrics(ctx context.Context, workspaceID string) (*model.StorageSnapshot, error)
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

func (s *defaultService) RecordSnapshot(ctx context.Context, snapshot *model.StorageSnapshot) error {
	return s.repo.UpsertSnapshot(ctx, snapshot)
}

func (s *defaultService) GetHistory(ctx context.Context, workspaceID string, from, to time.Time) ([]model.StorageSnapshot, error) {
	return s.repo.GetSnapshots(ctx, workspaceID, from, to)
}

func (s *defaultService) GetLatestMetrics(ctx context.Context, workspaceID string) (*model.StorageSnapshot, error) {
	return s.repo.GetLatestSnapshot(ctx, workspaceID)
}
