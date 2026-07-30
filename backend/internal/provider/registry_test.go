package provider

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/storage"
)

// MockProvider implements storage.StorageProvider for testing purposes.
type MockProvider struct{}

func (m *MockProvider) Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *storage.ObjectMetadata) (*storage.StorageResponse, error) {
	return &storage.StorageResponse{Bucket: bucket, Key: key, Size: size}, nil
}
func (m *MockProvider) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *MockProvider) Delete(ctx context.Context, bucket, key string) error { return nil }
func (m *MockProvider) Exists(ctx context.Context, bucket, key string) (bool, error) {
	return true, nil
}
func (m *MockProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]*storage.Object, error) {
	return nil, nil
}
func (m *MockProvider) CreateBucket(ctx context.Context, bucket string) error { return nil }
func (m *MockProvider) DeleteBucket(ctx context.Context, bucket string) error { return nil }
func (m *MockProvider) ListBuckets(ctx context.Context) ([]*storage.Bucket, error) {
	return []*storage.Bucket{{Name: "mock-bucket"}}, nil
}
func (m *MockProvider) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "http://mock-url", nil
}
func (m *MockProvider) CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	return nil, nil
}
func (m *MockProvider) MoveObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	return nil, nil
}
func (m *MockProvider) GetObjectMetadata(ctx context.Context, bucket, key string) (*storage.ObjectMetadata, error) {
	return nil, nil
}
func (m *MockProvider) SetObjectMetadata(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) error {
	return nil
}
func (m *MockProvider) GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error) {
	return nil, nil
}
func (m *MockProvider) SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error {
	return nil
}
func (m *MockProvider) CreateMultipartUpload(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) (*storage.MultipartUpload, error) {
	return nil, nil
}
func (m *MockProvider) UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*storage.UploadPart, error) {
	return nil, nil
}
func (m *MockProvider) CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*storage.UploadPart) (*storage.StorageResponse, error) {
	return nil, nil
}
func (m *MockProvider) AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error {
	return nil
}
func (m *MockProvider) ListObjectVersions(ctx context.Context, bucket, key string) ([]*storage.Object, error) {
	return nil, nil
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	mock := &MockProvider{}

	err := reg.Register("provider-1", mock)
	if err != nil {
		t.Fatalf("expected nil error on register, got %v", err)
	}

	got, err := reg.Get("provider-1")
	if err != nil {
		t.Fatalf("expected nil error on get, got %v", err)
	}
	if got != mock {
		t.Fatalf("expected returned provider to match registered instance")
	}
}

func TestRegistry_DuplicateRegister(t *testing.T) {
	reg := NewRegistry()
	mock := &MockProvider{}

	_ = reg.Register("provider-1", mock)
	err := reg.Register("provider-1", mock)
	if err == nil {
		t.Fatalf("expected error on duplicate register, got nil")
	}
}

func TestRegistry_GetNotFound(t *testing.T) {
	reg := NewRegistry()
	_, err := reg.Get("non-existent")
	if err == nil {
		t.Fatalf("expected error on missing provider, got nil")
	}
}

func TestRegistry_ConcurrentSafety(t *testing.T) {
	reg := NewRegistry()
	mock := &MockProvider{}

	done := make(chan bool)
	for i := 0; i < 50; i++ {
		go func(id int) {
			_ = reg.Register(string(rune(id)), mock)
			_, _ = reg.Get(string(rune(id)))
			_ = reg.List()
			done <- true
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}
}
