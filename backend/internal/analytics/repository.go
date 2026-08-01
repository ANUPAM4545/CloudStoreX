package analytics

import (
	"context"
	"time"

	"github.com/cloudstorex/backend/internal/analytics/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	UpsertSnapshot(ctx context.Context, snapshot *model.StorageSnapshot) error
	GetSnapshots(ctx context.Context, workspaceID string, from, to time.Time) ([]model.StorageSnapshot, error)
	GetLatestSnapshot(ctx context.Context, workspaceID string) (*model.StorageSnapshot, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) UpsertSnapshot(ctx context.Context, snapshot *model.StorageSnapshot) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "workspace_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"total_bytes", "total_objects", "bucket_count", "updated_at"}),
	}).Create(snapshot).Error
}

func (r *postgresRepository) GetSnapshots(ctx context.Context, workspaceID string, from, to time.Time) ([]model.StorageSnapshot, error) {
	var snapshots []model.StorageSnapshot
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND date >= ? AND date <= ?", workspaceID, from, to).
		Order("date ASC").
		Find(&snapshots).Error
	return snapshots, err
}

func (r *postgresRepository) GetLatestSnapshot(ctx context.Context, workspaceID string) (*model.StorageSnapshot, error) {
	var snapshot model.StorageSnapshot
	err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("date DESC").
		First(&snapshot).Error
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}
