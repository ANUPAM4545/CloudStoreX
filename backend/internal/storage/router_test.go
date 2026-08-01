package storage

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// mockRegistry implements ProviderRegistry for testing.
type mockRegistry struct {
	providers map[string]StorageProvider
}

func (m *mockRegistry) Get(id string) (StorageProvider, error) {
	prov, ok := m.providers[id]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return prov, nil
}

// mockProvider implements StorageProvider for testing.
type mockProvider struct {
	uploadFunc   func(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error)
	downloadFunc func(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	deleteFunc   func(ctx context.Context, bucket, key string) error
	existsFunc   func(ctx context.Context, bucket, key string) (bool, error)
}

func (m *mockProvider) Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, bucket, key, reader, size, meta)
	}
	return &StorageResponse{Bucket: bucket, Key: key, Size: size}, nil
}
func (m *mockProvider) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	if m.downloadFunc != nil {
		return m.downloadFunc(ctx, bucket, key)
	}
	return io.NopCloser(strings.NewReader("test data")), nil
}
func (m *mockProvider) Delete(ctx context.Context, bucket, key string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, bucket, key)
	}
	return nil
}
func (m *mockProvider) Exists(ctx context.Context, bucket, key string) (bool, error) {
	if m.existsFunc != nil {
		return m.existsFunc(ctx, bucket, key)
	}
	return true, nil
}
func (m *mockProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error) {
	return []*Object{{Key: "obj1", Bucket: bucket, Size: 100}}, nil
}
func (m *mockProvider) CreateBucket(ctx context.Context, bucket string) error { return nil }
func (m *mockProvider) DeleteBucket(ctx context.Context, bucket string) error { return nil }
func (m *mockProvider) ListBuckets(ctx context.Context) ([]*Bucket, error) {
	return []*Bucket{{Name: "test-bucket"}}, nil
}
func (m *mockProvider) GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "http://presigned-upload", nil
}

func (m *mockProvider) GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "http://presigned-download", nil
}
func (m *mockProvider) CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*StorageResponse, error) {
	return nil, nil
}
func (m *mockProvider) MoveObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*StorageResponse, error) {
	return nil, nil
}
func (m *mockProvider) GetObjectMetadata(ctx context.Context, bucket, key string) (*ObjectMetadata, error) {
	return nil, nil
}
func (m *mockProvider) SetObjectMetadata(ctx context.Context, bucket, key string, meta *ObjectMetadata) error {
	return nil
}
func (m *mockProvider) GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error) {
	return nil, nil
}
func (m *mockProvider) SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error {
	return nil
}
func (m *mockProvider) CreateMultipartUpload(ctx context.Context, bucket, key string, meta *ObjectMetadata) (*MultipartUpload, error) {
	return nil, nil
}
func (m *mockProvider) UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*UploadPart, error) {
	return nil, nil
}
func (m *mockProvider) CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*UploadPart) (*StorageResponse, error) {
	return nil, nil
}
func (m *mockProvider) AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error {
	return nil
}
func (m *mockProvider) ListObjectVersions(ctx context.Context, bucket, key string) ([]*Object, error) {
	return nil, nil
}

type mockPolicyEngine struct {
	id string
}

func (m mockPolicyEngine) Resolve(ctx context.Context, evalCtx *model.EvaluationContext, op string) (string, *model.RoutingDecision, error) {
	return m.id, nil, nil
}

func TestRouter_Upload(t *testing.T) {
	evaluator := mockPolicyEngine{id: "test-provider"}
	mockProv := &mockProvider{}
	registry := &mockRegistry{
		providers: map[string]StorageProvider{"test-provider": mockProv},
	}

	router := NewRouter(registry, evaluator)

	res, err := router.Upload(context.Background(), "my-bucket", "test.txt", strings.NewReader("hello"), 5, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.Bucket != "my-bucket" || res.Key != "test.txt" {
		t.Errorf("unexpected response: %+v", res)
	}
}

func TestRouter_ProviderNotFound(t *testing.T) {
	evaluator := mockPolicyEngine{id: "missing-provider"}
	registry := &mockRegistry{
		providers: map[string]StorageProvider{},
	}

	router := NewRouter(registry, evaluator)

	_, err := router.Upload(context.Background(), "my-bucket", "test.txt", strings.NewReader("hello"), 5, nil)
	if !errors.Is(err, ErrProviderNotFound) {
		t.Fatalf("expected ErrProviderNotFound, got %v", err)
	}
}
