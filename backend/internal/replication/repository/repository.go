package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/replication/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, repl *model.ObjectReplication) error
	Update(ctx context.Context, repl *model.ObjectReplication) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.ObjectReplication, error)
	ListByObjectID(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error)
	ListByStatus(ctx context.Context, status model.ReplicationStatus, limit int) ([]model.ObjectReplication, error)
}

type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new GORM-backed replication repository.
func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, repl *model.ObjectReplication) error {
	return r.db.WithContext(ctx).Create(repl).Error
}

func (r *postgresRepository) Update(ctx context.Context, repl *model.ObjectReplication) error {
	return r.db.WithContext(ctx).Save(repl).Error
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.ObjectReplication, error) {
	var repl model.ObjectReplication
	if err := r.db.WithContext(ctx).First(&repl, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &repl, nil
}

func (r *postgresRepository) ListByObjectID(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error) {
	var repls []model.ObjectReplication
	err := r.db.WithContext(ctx).Where("object_id = ?", objectID).Order("created_at DESC").Find(&repls).Error
	return repls, err
}

func (r *postgresRepository) ListByStatus(ctx context.Context, status model.ReplicationStatus, limit int) ([]model.ObjectReplication, error) {
	var repls []model.ObjectReplication
	err := r.db.WithContext(ctx).Where("status = ?", status).Limit(limit).Find(&repls).Error
	return repls, err
}
