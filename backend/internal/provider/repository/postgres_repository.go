package repository

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/provider/model"
	"github.com/cloudstorex/backend/internal/storage"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, provider *model.Provider) error {
	return r.db.WithContext(ctx).Create(provider).Error
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*model.Provider, error) {
	var provider model.Provider
	if err := r.db.WithContext(ctx).First(&provider, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrProviderNotFound
		}
		return nil, err
	}
	return &provider, nil
}

func (r *postgresRepository) GetByName(ctx context.Context, workspaceID, name string) (*model.Provider, error) {
	var provider model.Provider
	if err := r.db.WithContext(ctx).First(&provider, "workspace_id = ? AND provider_name = ?", workspaceID, name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrProviderNotFound
		}
		return nil, err
	}
	return &provider, nil
}

func (r *postgresRepository) List(ctx context.Context, workspaceID string) ([]*model.Provider, error) {
	var providers []*model.Provider
	if err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *postgresRepository) GetDefault(ctx context.Context, workspaceID string) (*model.Provider, error) {
	var provider model.Provider
	if err := r.db.WithContext(ctx).First(&provider, "workspace_id = ? AND is_default = ?", workspaceID, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrProviderNotFound
		}
		return nil, err
	}
	return &provider, nil
}

func (r *postgresRepository) Update(ctx context.Context, provider *model.Provider) error {
	return r.db.WithContext(ctx).Save(provider).Error
}

func (r *postgresRepository) UpdateHealth(ctx context.Context, id string, status model.ProviderStatus, health model.ProviderHealth, latencyMs int64, lastError string) error {
	updates := map[string]interface{}{
		"status":            status,
		"health":            health,
		"latency_ms":        latencyMs,
		"last_error":        lastError,
		"last_health_check": gorm.Expr("CURRENT_TIMESTAMP"),
	}
	res := r.db.WithContext(ctx).Model(&model.Provider{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return storage.ErrProviderNotFound
	}
	return nil
}

func (r *postgresRepository) SetDefault(ctx context.Context, workspaceID, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset current default
		if err := tx.Model(&model.Provider{}).Where("workspace_id = ? AND is_default = ?", workspaceID, true).Update("is_default", false).Error; err != nil {
			return err
		}
		// Set new default
		res := tx.Model(&model.Provider{}).Where("workspace_id = ? AND id = ?", workspaceID, id).Update("is_default", true)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return storage.ErrProviderNotFound
		}
		return nil
	})
}

func (r *postgresRepository) Enable(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Model(&model.Provider{}).Where("id = ?", id).Update("is_enabled", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return storage.ErrProviderNotFound
	}
	return nil
}

func (r *postgresRepository) Disable(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Model(&model.Provider{}).Where("id = ?", id).Update("is_enabled", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return storage.ErrProviderNotFound
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&model.Provider{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return storage.ErrProviderNotFound
	}
	return nil
}
