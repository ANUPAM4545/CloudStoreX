package service

import (
	"context"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/model"
	"github.com/cloudstorex/backend/internal/metadata/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type benchMockMetadataRepo struct{}

func (m *benchMockMetadataRepo) CreateObject(ctx context.Context, obj *model.Object) error { return nil }
func (m *benchMockMetadataRepo) CreateObjectVersion(ctx context.Context, version *model.ObjectVersion) error {
	return nil
}
func (m *benchMockMetadataRepo) UpdateObject(ctx context.Context, obj *model.Object) error { return nil }
func (m *benchMockMetadataRepo) DeleteObject(ctx context.Context, id string) error         { return nil }
func (m *benchMockMetadataRepo) FindObjectByKey(ctx context.Context, bucketID, objectKey string) (*model.Object, error) {
	return nil, nil
}
func (m *benchMockMetadataRepo) FindObjectByID(ctx context.Context, id string) (*model.Object, error) {
	return nil, nil
}
func (m *benchMockMetadataRepo) ListObjectVersions(ctx context.Context, id string) ([]model.ObjectVersion, error) {
	return nil, nil
}
func (m *benchMockMetadataRepo) SearchObjects(ctx context.Context, query dto.SearchQuery) ([]model.Object, int64, error) {
	return []model.Object{{ID: uuid.New(), ObjectKey: "search-1"}, {ID: uuid.New(), ObjectKey: "search-2"}}, 2, nil
}
func (m *benchMockMetadataRepo) FindObjectsForExpiration(ctx context.Context, bucketID string, prefix string, olderThan time.Time) ([]model.Object, error) {
	return nil, nil
}
func (m *benchMockMetadataRepo) CreateBucket(ctx context.Context, bucket *model.Bucket) error {
	return nil
}
func (m *benchMockMetadataRepo) FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error) {
	wsUUID, _ := uuid.Parse(workspaceID)
	return &model.Bucket{
		ID:          uuid.New(),
		WorkspaceID: wsUUID,
		Name:        bucketName,
	}, nil
}
func (m *benchMockMetadataRepo) ListBuckets(ctx context.Context, workspaceID string) ([]model.Bucket, error) {
	return nil, nil
}
func (m *benchMockMetadataRepo) WithTx(tx *gorm.DB) repository.MetadataRepository { return m }
func (m *benchMockMetadataRepo) DB() *gorm.DB                                     { return nil }

// BenchmarkMetadataService_FindBucketByName benchmarks bucket lookup latency and memory allocations.
func BenchmarkMetadataService_FindBucketByName(b *testing.B) {
	svc := NewMetadataService(&benchMockMetadataRepo{}, nil, nil)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = svc.FindBucketByName(ctx, "00000000-0000-0000-0000-000000000100", "bench-catalog-bucket")
	}
}

// BenchmarkMetadataService_FindBucketByNameParallel benchmarks parallel reader scalability for catalog lookup.
func BenchmarkMetadataService_FindBucketByNameParallel(b *testing.B) {
	svc := NewMetadataService(&benchMockMetadataRepo{}, nil, nil)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = svc.FindBucketByName(ctx, "00000000-0000-0000-0000-000000000100", "bench-catalog-bucket")
		}
	})
}

// BenchmarkMetadataService_SearchObjects benchmarks metadata search evaluation.
func BenchmarkMetadataService_SearchObjects(b *testing.B) {
	svc := NewMetadataService(&benchMockMetadataRepo{}, nil, nil)
	ctx := context.Background()

	query := dto.SearchQuery{
		WorkspaceID: "00000000-0000-0000-0000-000000000100",
		Prefix:      "search-item",
		Limit:       20,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = svc.SearchObjects(ctx, query)
	}
}
