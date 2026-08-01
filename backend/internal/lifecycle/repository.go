package lifecycle

import (
	"context"

	"github.com/cloudstorex/backend/internal/lifecycle/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateRule(ctx context.Context, rule *model.LifecycleRule) error
	ListRulesByBucket(ctx context.Context, bucketID uuid.UUID) ([]model.LifecycleRule, error)
	GetActiveRules(ctx context.Context) ([]model.LifecycleRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateRule(ctx context.Context, rule *model.LifecycleRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *postgresRepository) ListRulesByBucket(ctx context.Context, bucketID uuid.UUID) ([]model.LifecycleRule, error) {
	var rules []model.LifecycleRule
	err := r.db.WithContext(ctx).Where("bucket_id = ?", bucketID).Find(&rules).Error
	return rules, err
}

func (r *postgresRepository) GetActiveRules(ctx context.Context) ([]model.LifecycleRule, error) {
	var rules []model.LifecycleRule
	err := r.db.WithContext(ctx).Where("status = ?", "ACTIVE").Find(&rules).Error
	return rules, err
}

func (r *postgresRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.LifecycleRule{}, "id = ?", id).Error
}
