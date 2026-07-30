package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"


	"github.com/cloudstorex/backend/internal/policy"
	"github.com/cloudstorex/backend/internal/provider"
	"github.com/cloudstorex/backend/internal/storage"
	minio "github.com/minio/minio-go/v7"

)

// mockMinIOClient implements MinIOClient for unit testing.
type mockMinIOClient struct {
	putFunc       func(ctx context.Context, bucket, key string, r io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	getFunc       func(ctx context.Context, bucket, key string, opts minio.GetObjectOptions) (*minio.Object, error)
	removeFunc    func(ctx context.Context, bucket, key string, opts minio.RemoveObjectOptions) error
	statFunc      func(ctx context.Context, bucket, key string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	listFunc      func(ctx context.Context, bucket string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	makeBucket      func(ctx context.Context, bucket string, opts minio.MakeBucketOptions) error
	bucketExists    func(ctx context.Context, bucket string) (bool, error)
	listBucketsFunc func(ctx context.Context) ([]minio.BucketInfo, error)
	objects         map[string]map[string][]byte // bucket -> key -> data
}

func newMockMinIOClient() *mockMinIOClient {
	return &mockMinIOClient{
		objects: make(map[string]map[string][]byte),
	}
}

func (m *mockMinIOClient) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	if m.listBucketsFunc != nil {
		return m.listBucketsFunc(ctx)
	}
	var res []minio.BucketInfo
	for k := range m.objects {
		res = append(res, minio.BucketInfo{Name: k})
	}
	return res, nil
}

func (m *mockMinIOClient) BucketExists(ctx context.Context, bucket string) (bool, error) {
	if m.bucketExists != nil {
		return m.bucketExists(ctx, bucket)
	}
	_, ok := m.objects[bucket]
	return ok, nil
}

func (m *mockMinIOClient) MakeBucket(ctx context.Context, bucket string, opts minio.MakeBucketOptions) error {
	if m.makeBucket != nil {
		return m.makeBucket(ctx, bucket, opts)
	}
	if _, ok := m.objects[bucket]; ok {
		return minio.ErrorResponse{Code: "BucketAlreadyExists"}
	}
	m.objects[bucket] = make(map[string][]byte)
	return nil
}

func (m *mockMinIOClient) PutObject(ctx context.Context, bucket, key string, r io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	if m.putFunc != nil {
		return m.putFunc(ctx, bucket, key, r, size, opts)
	}
	data, _ := io.ReadAll(r)
	if m.objects[bucket] == nil {
		m.objects[bucket] = make(map[string][]byte)
	}
	m.objects[bucket][key] = data
	return minio.UploadInfo{Bucket: bucket, Key: key, Size: int64(len(data)), ETag: "mock-etag"}, nil
}

func (m *mockMinIOClient) GetObject(ctx context.Context, bucket, key string, opts minio.GetObjectOptions) (*minio.Object, error) {
	// For simplicity in mock without real socket, we test download via StatObject check
	return nil, nil
}

func (m *mockMinIOClient) RemoveObject(ctx context.Context, bucket, key string, opts minio.RemoveObjectOptions) error {
	if m.removeFunc != nil {
		return m.removeFunc(ctx, bucket, key, opts)
	}
	if _, ok := m.objects[bucket]; ok {
		delete(m.objects[bucket], key)
	}
	return nil
}

func (m *mockMinIOClient) StatObject(ctx context.Context, bucket, key string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	if m.statFunc != nil {
		return m.statFunc(ctx, bucket, key, opts)
	}
	b, ok := m.objects[bucket]
	if !ok {
		return minio.ObjectInfo{}, minio.ErrorResponse{Code: "NoSuchBucket"}
	}
	data, ok := b[key]
	if !ok {
		return minio.ObjectInfo{}, minio.ErrorResponse{Code: "NoSuchKey"}
	}
	return minio.ObjectInfo{Key: key, Size: int64(len(data))}, nil
}

func (m *mockMinIOClient) ListObjects(ctx context.Context, bucket string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	ch := make(chan minio.ObjectInfo, 10)
	go func() {
		defer close(ch)
		if b, ok := m.objects[bucket]; ok {
			for k, v := range b {
				if strings.HasPrefix(k, opts.Prefix) {
					ch <- minio.ObjectInfo{Key: k, Size: int64(len(v))}
				}
			}
		}
	}()
	return ch
}

func TestMinIOProvider_InitAndVerify(t *testing.T) {
	mockClient := newMockMinIOClient()
	cfg := &Config{DefaultBucket: "cloudstorex-test"}

	prov, err := NewProviderWithClient(context.Background(), mockClient, cfg, nil)
	if err != nil {
		t.Fatalf("expected nil error on NewProviderWithClient, got %v", err)
	}
	if prov == nil {
		t.Fatalf("expected provider instance, got nil")
	}

	// Verify default bucket was created
	exists, err := mockClient.BucketExists(context.Background(), "cloudstorex-test")
	if err != nil || !exists {
		t.Fatalf("expected default bucket to be created automatically")
	}
}

func TestMinIOProvider_ErrorTranslation(t *testing.T) {
	mockClient := newMockMinIOClient()
	cfg := &Config{DefaultBucket: "cloudstorex-test"}

	prov, _ := NewProviderWithClient(context.Background(), mockClient, cfg, nil)

	// Stat missing object -> ErrObjectNotFound
	exists, err := prov.Exists(context.Background(), "cloudstorex-test", "missing.txt")
	if err != nil {
		t.Fatalf("expected nil error when checking missing object existence, got %v", err)
	}
	if exists {
		t.Fatalf("expected exists=false for missing object")
	}

	// Simulate upload failure
	mockClient.putFunc = func(ctx context.Context, b, k string, r io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
		return minio.UploadInfo{}, minio.ErrorResponse{Code: "InternalError", StatusCode: http.StatusInternalServerError}
	}

	_, err = prov.Upload(context.Background(), "cloudstorex-test", "test.txt", strings.NewReader("data"), 4, nil)
	if !errors.Is(err, storage.ErrUploadFailed) {
		t.Fatalf("expected ErrUploadFailed, got %v", err)
	}
}

func TestMinIOProvider_EndToEndPipelineIntegration(t *testing.T) {
	// 1. Setup MinIO Provider with Mock MinIOClient
	mockClient := newMockMinIOClient()
	cfg := &Config{DefaultBucket: "cloudstorex-default"}
	minioProv, err := NewProviderWithClient(context.Background(), mockClient, cfg, nil)
	if err != nil {
		t.Fatalf("failed to init minio provider: %v", err)
	}

	// 2. Register MinIO with Provider Registry
	registry := provider.NewRegistry()
	err = registry.Register("minio", minioProv)
	if err != nil {
		t.Fatalf("failed to register minio provider: %v", err)
	}

	// 3. Configure Policy Evaluator to route to "minio"
	evaluator := policy.NewDefaultEvaluator("minio")

	// 4. Instantiate central Storage Router
	router := storage.NewRouter(registry, evaluator)

	// 5. Instantiate high-level Storage Service
	service := storage.NewService(router)

	ctx := context.Background()

	// Step A: Upload object through Storage Service
	data := []byte("CloudStoreX minio integration test content")
	res, err := service.UploadObject(ctx, "cloudstorex-default", "docs/architecture.txt", bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		t.Fatalf("expected nil error on UploadObject, got %v", err)
	}
	if res.Bucket != "cloudstorex-default" || res.Key != "docs/architecture.txt" || res.Size != int64(len(data)) {
		t.Errorf("unexpected upload response: %+v", res)
	}

	// Step B: Verify Object Exists through Storage Service
	exists, err := service.ObjectExists(ctx, "cloudstorex-default", "docs/architecture.txt")
	if err != nil {
		t.Fatalf("expected nil error on ObjectExists, got %v", err)
	}
	if !exists {
		t.Fatalf("expected object to exist after upload")
	}

	// Step C: List Objects through Storage Service
	objs, err := service.ListObjects(ctx, "cloudstorex-default", "docs/")
	if err != nil {
		t.Fatalf("expected nil error on ListObjects, got %v", err)
	}
	if len(objs) != 1 || objs[0].Key != "docs/architecture.txt" {
		t.Errorf("unexpected list result: %+v", objs)
	}

	// Step D: Delete Object through Storage Service
	err = service.DeleteObject(ctx, "cloudstorex-default", "docs/architecture.txt")
	if err != nil {
		t.Fatalf("expected nil error on DeleteObject, got %v", err)
	}

	// Step E: Verify Object is Gone
	exists, _ = service.ObjectExists(ctx, "cloudstorex-default", "docs/architecture.txt")
	if exists {
		t.Fatalf("expected object to be gone after delete")
	}
}
