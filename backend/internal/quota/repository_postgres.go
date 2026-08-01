package quota

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/quota/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetWorkspaceQuota(ctx context.Context, workspaceID string) (*model.WorkspaceQuota, error) {
	var q model.WorkspaceQuota
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).First(&q).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("quota not found")
		}
		return nil, err
	}
	return &q, nil
}

func (r *postgresRepository) UpsertWorkspaceQuota(ctx context.Context, quota *model.WorkspaceQuota) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "workspace_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"max_bytes", "max_objects", "updated_at"}),
	}).Create(quota).Error
}

func (r *postgresRepository) IncrementUsage(ctx context.Context, workspaceID string, bytes int64, objects int64) error {
	// Increment using GORM expressions to avoid race conditions
	res := r.db.WithContext(ctx).Model(&model.WorkspaceQuota{}).
		Where("workspace_id = ?", workspaceID).
		Updates(map[string]interface{}{
			"bytes_used":   gorm.Expr("bytes_used + ?", bytes),
			"object_count": gorm.Expr("object_count + ?", objects),
		})
		
	if res.Error != nil {
		return res.Error
	}
	
	if res.RowsAffected == 0 {
		// Create the quota record if it doesn't exist
		quota := &model.WorkspaceQuota{
			WorkspaceID: workspaceID,
			BytesUsed:   bytes,
			ObjectCount: objects,
		}
		return r.UpsertWorkspaceQuota(ctx, quota)
	}
	
	return nil
}

func (r *postgresRepository) DecrementUsage(ctx context.Context, workspaceID string, bytes int64, objects int64) error {
	return r.db.WithContext(ctx).Model(&model.WorkspaceQuota{}).
		Where("workspace_id = ?", workspaceID).
		Updates(map[string]interface{}{
			"bytes_used":   gorm.Expr("GREATEST(0, bytes_used - ?)", bytes),
			"object_count": gorm.Expr("GREATEST(0, object_count - ?)", objects),
		}).Error
}
