package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/chaos/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateExperiment(ctx context.Context, exp *model.ChaosExperiment) error
	UpdateExperiment(ctx context.Context, exp *model.ChaosExperiment) error
	GetExperimentByID(ctx context.Context, id uuid.UUID) (*model.ChaosExperiment, error)
	ListExperimentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.ChaosExperiment, error)
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new GORM-backed chaos engineering repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateExperiment(ctx context.Context, exp *model.ChaosExperiment) error {
	return r.db.WithContext(ctx).Create(exp).Error
}

func (r *postgresRepository) UpdateExperiment(ctx context.Context, exp *model.ChaosExperiment) error {
	return r.db.WithContext(ctx).Save(exp).Error
}

func (r *postgresRepository) GetExperimentByID(ctx context.Context, id uuid.UUID) (*model.ChaosExperiment, error) {
	var exp model.ChaosExperiment
	if err := r.db.WithContext(ctx).First(&exp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *postgresRepository) ListExperimentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.ChaosExperiment, error) {
	var res []model.ChaosExperiment
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&res).Error
	return res, err
}
