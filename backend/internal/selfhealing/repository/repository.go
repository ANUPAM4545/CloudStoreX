package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/selfhealing/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateJob(ctx context.Context, job *model.RepairJob) error
	UpdateJob(ctx context.Context, job *model.RepairJob) error
	GetJobByID(ctx context.Context, id uuid.UUID) (*model.RepairJob, error)
	ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error)
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new GORM-backed self-healing repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateJob(ctx context.Context, job *model.RepairJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *postgresRepository) UpdateJob(ctx context.Context, job *model.RepairJob) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *postgresRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*model.RepairJob, error) {
	var job model.RepairJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *postgresRepository) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error) {
	var jobs []model.RepairJob
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}
