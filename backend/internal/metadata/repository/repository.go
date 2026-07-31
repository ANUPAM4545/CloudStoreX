package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/model"
	"gorm.io/gorm"
)

// MetadataRepository handles interactions with the Postgres metadata catalog.
type MetadataRepository interface {
	CreateObject(ctx context.Context, obj *model.Object) error
	UpdateObject(ctx context.Context, obj *model.Object) error
	DeleteObject(ctx context.Context, id string) error // Soft delete
	FindObjectByKey(ctx context.Context, bucketID, objectKey string) (*model.Object, error)
	FindObjectByID(ctx context.Context, id string) (*model.Object, error)
	SearchObjects(ctx context.Context, query dto.SearchQuery) ([]model.Object, int64, error)
	
	CreateBucket(ctx context.Context, bucket *model.Bucket) error
	FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error)
	ListBuckets(ctx context.Context, workspaceID string) ([]model.Bucket, error)
	
	// Transaction management
	WithTx(tx *gorm.DB) MetadataRepository
	DB() *gorm.DB
}

type postgresMetadataRepository struct {
	db *gorm.DB
}

func NewPostgresMetadataRepository(db *gorm.DB) MetadataRepository {
	return &postgresMetadataRepository{db: db}
}

func (r *postgresMetadataRepository) WithTx(tx *gorm.DB) MetadataRepository {
	return &postgresMetadataRepository{db: tx}
}

func (r *postgresMetadataRepository) DB() *gorm.DB {
	return r.db
}

func (r *postgresMetadataRepository) CreateObject(ctx context.Context, obj *model.Object) error {
	return r.db.WithContext(ctx).Create(obj).Error
}

func (r *postgresMetadataRepository) UpdateObject(ctx context.Context, obj *model.Object) error {
	return r.db.WithContext(ctx).Save(obj).Error
}

func (r *postgresMetadataRepository) DeleteObject(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.Object{}).Where("id = ?", id).Update("is_deleted", true).Delete(&model.Object{}, "id = ?", id).Error
}

func (r *postgresMetadataRepository) FindObjectByKey(ctx context.Context, bucketID, objectKey string) (*model.Object, error) {
	var obj model.Object
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Metadata").
		Where("bucket_id = ? AND object_key = ? AND is_deleted = false", bucketID, objectKey).
		First(&obj).Error
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (r *postgresMetadataRepository) FindObjectByID(ctx context.Context, id string) (*model.Object, error) {
	var obj model.Object
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Preload("Metadata").
		Where("id = ? AND is_deleted = false", id).
		First(&obj).Error
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (r *postgresMetadataRepository) SearchObjects(ctx context.Context, query dto.SearchQuery) ([]model.Object, int64, error) {
	db := r.db.WithContext(ctx).Model(&model.Object{}).Where("is_deleted = false")

	if query.WorkspaceID != "" {
		// Join bucket to filter by workspace
		db = db.Joins("JOIN buckets ON buckets.id = objects.bucket_id").
			Where("buckets.workspace_id = ?", query.WorkspaceID)
	}

	if query.BucketID != "" {
		db = db.Where("objects.bucket_id = ?", query.BucketID)
	}
	if query.Prefix != "" {
		db = db.Where("objects.object_key LIKE ?", query.Prefix+"%")
	}
	if query.MimeType != "" {
		db = db.Where("objects.mime_type = ?", query.MimeType)
	}
	if query.OwnerID != "" {
		db = db.Where("objects.owner_id = ?", query.OwnerID)
	}
	if query.ProviderID != "" {
		db = db.Where("objects.provider_id = ?", query.ProviderID)
	}
	if query.StorageClass != "" {
		db = db.Where("objects.storage_class = ?", query.StorageClass)
	}
	if query.Status != "" {
		db = db.Where("objects.status = ?", query.Status)
	}

	if len(query.Tags) > 0 {
		for k, v := range query.Tags {
			db = db.Where("EXISTS (SELECT 1 FROM object_tags WHERE object_tags.object_id = objects.id AND object_tags.key = ? AND object_tags.value = ?)", k, v)
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Offset > 0 {
		db = db.Offset(query.Offset)
	}

	var results []model.Object
	err := db.Preload("Tags").Preload("Metadata").Order("objects.created_at DESC").Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *postgresMetadataRepository) CreateBucket(ctx context.Context, bucket *model.Bucket) error {
	return r.db.WithContext(ctx).Create(bucket).Error
}

func (r *postgresMetadataRepository) FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error) {
	var bucket model.Bucket
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND name = ?", workspaceID, bucketName).First(&bucket).Error
	if err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (r *postgresMetadataRepository) ListBuckets(ctx context.Context, workspaceID string) ([]model.Bucket, error) {
	var buckets []model.Bucket
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("created_at DESC").Find(&buckets).Error
	return buckets, err
}
