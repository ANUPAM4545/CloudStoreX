package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/disasterrecovery/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePlan(ctx context.Context, plan *model.RecoveryPlan) error
	GetPlanByID(ctx context.Context, id uuid.UUID) (*model.RecoveryPlan, error)
	ListPlansByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.RecoveryPlan, error)

	CreateJob(ctx context.Context, job *model.RecoveryJob) error
	UpdateJob(ctx context.Context, job *model.RecoveryJob) error
	GetJobByID(ctx context.Context, id uuid.UUID) (*model.RecoveryJob, error)
	ListJobsByPlan(ctx context.Context, planID uuid.UUID) ([]model.RecoveryJob, error)
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new GORM-backed disaster recovery repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreatePlan(ctx context.Context, plan *model.RecoveryPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *postgresRepository) GetPlanByID(ctx context.Context, id uuid.UUID) (*model.RecoveryPlan, error) {
	var plan model.RecoveryPlan
	if err := r.db.WithContext(ctx).First(&plan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *postgresRepository) ListPlansByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.RecoveryPlan, error) {
	var plans []model.RecoveryPlan
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&plans).Error
	return plans, err
}

func (r *postgresRepository) CreateJob(ctx context.Context, job *model.RecoveryJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *postgresRepository) UpdateJob(ctx context.Context, job *model.RecoveryJob) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *postgresRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*model.RecoveryJob, error) {
	var job model.RecoveryJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *postgresRepository) ListJobsByPlan(ctx context.Context, planID uuid.UUID) ([]model.RecoveryJob, error) {
	var jobs []model.RecoveryJob
	err := r.db.WithContext(ctx).Where("plan_id = ?", planID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}
