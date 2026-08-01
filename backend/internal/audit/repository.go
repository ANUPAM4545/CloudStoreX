package audit

import (
	"context"

	"github.com/cloudstorex/backend/internal/audit/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	ListByWorkspace(ctx context.Context, workspaceID string, limit, offset int) ([]model.AuditLog, int64, error)
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *postgresRepository) ListByWorkspace(ctx context.Context, workspaceID string, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("workspace_id = ?", workspaceID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}
