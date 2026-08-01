package lifecycle

import (
	"context"
	"log/slog"

	"github.com/cloudstorex/backend/internal/lifecycle/model"
	"github.com/google/uuid"
)

type Service interface {
	CreateRule(ctx context.Context, rule *model.LifecycleRule) error
	ListRules(ctx context.Context, bucketID uuid.UUID) ([]model.LifecycleRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
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

func (s *defaultService) CreateRule(ctx context.Context, rule *model.LifecycleRule) error {
	return s.repo.CreateRule(ctx, rule)
}

func (s *defaultService) ListRules(ctx context.Context, bucketID uuid.UUID) ([]model.LifecycleRule, error) {
	return s.repo.ListRulesByBucket(ctx, bucketID)
}

func (s *defaultService) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteRule(ctx, id)
}
