package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/backup/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateJob(ctx context.Context, job *model.BackupJob) error
	UpdateJob(ctx context.Context, job *model.BackupJob) error
	GetJobByID(ctx context.Context, id uuid.UUID) (*model.BackupJob, error)
	ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupJob, error)

	CreateRecord(ctx context.Context, rec *model.BackupRecord) error
	ListRecordsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupRecord, error)
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new GORM-backed backup repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateJob(ctx context.Context, job *model.BackupJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *postgresRepository) UpdateJob(ctx context.Context, job *model.BackupJob) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *postgresRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*model.BackupJob, error) {
	var job model.BackupJob
	if err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *postgresRepository) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupJob, error) {
	var jobs []model.BackupJob
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&jobs).Error
	return jobs, err
}

func (r *postgresRepository) CreateRecord(ctx context.Context, rec *model.BackupRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *postgresRepository) ListRecordsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupRecord, error) {
	var records []model.BackupRecord
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&records).Error
	return records, err
}
