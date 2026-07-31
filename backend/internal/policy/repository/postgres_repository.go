package repository

import (
	"context"
	"fmt"

	"github.com/cloudstorex/backend/internal/policy/model"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new Postgres backed policy repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, p *model.Policy) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*model.Policy, error) {
	var p model.Policy
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("policy not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresRepository) List(ctx context.Context, workspaceID string) ([]*model.Policy, error) {
	var policies []*model.Policy
	q := r.db.WithContext(ctx)
	if workspaceID != "" {
		q = q.Where("workspace_id = ?", workspaceID)
	}
	if err := q.Order("priority asc").Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}

func (r *postgresRepository) ListEnabled(ctx context.Context, workspaceID string) ([]*model.Policy, error) {
	var policies []*model.Policy
	q := r.db.WithContext(ctx).Where("enabled = ?", true)
	if workspaceID != "" {
		q = q.Where("workspace_id = ?", workspaceID)
	}
	if err := q.Order("priority asc").Find(&policies).Error; err != nil {
		return nil, err
	}
	return policies, nil
}

func (r *postgresRepository) Update(ctx context.Context, p *model.Policy) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Policy{}, "id = ?", id).Error
}

func (r *postgresRepository) CreateRoutingDecision(ctx context.Context, d *model.RoutingDecision) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *postgresRepository) ListRoutingDecisions(ctx context.Context, workspaceID string, limit, offset int) ([]*model.RoutingDecision, int64, error) {
	var decisions []*model.RoutingDecision
	var total int64
	
	q := r.db.WithContext(ctx).Model(&model.RoutingDecision{})
	if workspaceID != "" {
		q = q.Where("workspace_id = ?", workspaceID)
	}
	
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	if err := q.Order("timestamp desc").Limit(limit).Offset(offset).Find(&decisions).Error; err != nil {
		return nil, 0, err
	}
	return decisions, total, nil
}
