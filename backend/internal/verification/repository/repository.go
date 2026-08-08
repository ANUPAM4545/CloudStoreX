package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/verification/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateJob(ctx context.Context, job *model.ConsistencyJob) error
	UpdateJob(ctx context.Context, job *model.ConsistencyJob) error
	GetJobByID(ctx context.Context, id uuid.UUID) (*model.ConsistencyJob, error)
	ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.ConsistencyJob, error)

	CreateDiscrepancy(ctx context.Context, rec *model.DiscrepancyRecord) error
	ListDiscrepanciesByWorkspace(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]model.DiscrepancyRecord, error)
	MarkDiscrepancyResolved(ctx context.Context, id uuid.UUID) error
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new GORM-backed consistency verification repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateJob(ctx context.Context, job *model.ConsistencyJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *postgresRepository) UpdateJob(ctx context.Context, job *model.ConsistencyJob) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *postgresRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*model.ConsistencyJob, error) {
	var job model.ConsistencyJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *postgresRepository) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.ConsistencyJob, error) {
	var jobs []model.ConsistencyJob
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *postgresRepository) CreateDiscrepancy(ctx context.Context, rec *model.DiscrepancyRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *postgresRepository) ListDiscrepanciesByWorkspace(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]model.DiscrepancyRecord, error) {
	var records []model.DiscrepancyRecord
	query := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID)
	if unresolvedOnly {
		query = query.Where("resolved = ?", false)
	}
	err := query.Order("created_at DESC").Find(&records).Error
	return records, err
}

func (r *postgresRepository) MarkDiscrepancyResolved(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.DiscrepancyRecord{}).Where("id = ?", id).Update("resolved", true).Error
}
