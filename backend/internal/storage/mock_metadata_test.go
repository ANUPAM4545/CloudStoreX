package storage

import (
	"context"
	"time"
	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/model"
	"github.com/google/uuid"
)

type mockMetadataService struct{}

func (m *mockMetadataService) ListObjectVersions(ctx context.Context, id string) ([]dto.ObjectVersionDTO, error) {
	return nil, nil
}

func (m *mockMetadataService) CreateObjectMetadata(ctx context.Context, workspaceID string, bucketID uuid.UUID, objectKey, providerObjectKey string, size int64, mimeType, etag, providerID string, ownerID *uuid.UUID, tags map[string]string, meta map[string]string) (*model.Object, error) {
	return &model.Object{ID: uuid.New(), ObjectKey: objectKey}, nil
}
func (m *mockMetadataService) UpdateObjectMetadata(ctx context.Context, id string, mimeType string, storageClass string, tags map[string]string, custom map[string]string) (*model.Object, error) {
	return nil, nil
}
func (m *mockMetadataService) FindObjectByKey(ctx context.Context, bucketID, objectKey string) (*model.Object, error) {
	return &model.Object{ID: uuid.New(), ObjectKey: objectKey}, nil
}
func (m *mockMetadataService) FindObjectByID(ctx context.Context, id string) (*model.Object, error) {
	return nil, nil
}
func (m *mockMetadataService) GetObjectMetadata(ctx context.Context, id string) (map[string]string, error) {
	return nil, nil
}
func (m *mockMetadataService) SearchObjects(ctx context.Context, query dto.SearchQuery) ([]dto.ObjectMetaDTO, int64, error) {
	return []dto.ObjectMetaDTO{{ObjectKey: "test-key"}}, 1, nil
}
func (m *mockMetadataService) TagObject(ctx context.Context, id string, tags map[string]string) error {
	return nil
}
func (m *mockMetadataService) UntagObject(ctx context.Context, id string, keys []string) error {
	return nil
}
func (m *mockMetadataService) SoftDeleteObject(ctx context.Context, workspaceID, id string) error {
	return nil
}
func (m *mockMetadataService) RestoreObject(ctx context.Context, id string) error {
	return nil
}
func (m *mockMetadataService) FindObjectsForExpiration(ctx context.Context, bucketID string, prefix string, olderThan time.Time) ([]model.Object, error) {
	return nil, nil
}
func (m *mockMetadataService) CreateBucket(ctx context.Context, workspaceID uuid.UUID, providerID, name, region string) (*model.Bucket, error) {
	return nil, nil
}
func (m *mockMetadataService) FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error) {
	return &model.Bucket{ID: uuid.New(), Name: bucketName}, nil
}
func (m *mockMetadataService) DeleteBucket(ctx context.Context, workspaceID, bucketName string) error {
	return nil
}
func (m *mockMetadataService) ListBuckets(ctx context.Context, workspaceID string) ([]model.Bucket, error) {
	return nil, nil
}
