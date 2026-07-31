package storage

import (
	"context"
	"strings"
	"testing"
	"time"
)

// mockRouter implements Router for service testing.
type mockRouter struct {
	StorageProvider
}

func TestService_UploadAndDownloadObject(t *testing.T) {
	mockProv := &mockProvider{}
	router := &mockRouter{StorageProvider: mockProv}
	svc := NewService(router, &mockMetadataService{})

	ctx := context.WithValue(context.Background(), CtxKeyRequestID, "req-12345")
	ctx = context.WithValue(ctx, CtxKeyWorkspaceID, "ws-67890")

	res, err := svc.UploadObject(ctx, "bucket-test", "key-test", strings.NewReader("content"), 7, nil)
	if err != nil {
		t.Fatalf("expected nil error on UploadObject, got %v", err)
	}
	if res.Key != "key-test" {
		t.Errorf("unexpected key in response: %s", res.Key)
	}

	rc, err := svc.DownloadObject(ctx, "bucket-test", "key-test")
	if err != nil {
		t.Fatalf("expected nil error on DownloadObject, got %v", err)
	}
	if rc != nil {
		rc.Close()
	}
}

func TestService_AllOperations(t *testing.T) {
	mockProv := &mockProvider{}
	router := &mockRouter{StorageProvider: mockProv}
	svc := NewService(router, &mockMetadataService{})

	ctx := context.Background()

	if err := svc.CreateBucket(ctx, "new-bucket"); err != nil {
		t.Errorf("unexpected error on CreateBucket: %v", err)
	}
	if buckets, err := svc.ListBuckets(ctx); err != nil || len(buckets) == 0 {
		t.Errorf("unexpected error on ListBuckets: %v", err)
	}
	if exists, err := svc.ObjectExists(ctx, "new-bucket", "test-key"); err != nil || !exists {
		t.Errorf("unexpected error on ObjectExists: %v", err)
	}
	if objs, err := svc.ListObjects(ctx, "new-bucket", ""); err != nil || len(objs) == 0 {
		t.Errorf("unexpected error on ListObjects: %v", err)
	}
	if url, err := svc.GeneratePresignedURL(ctx, "new-bucket", "test-key", time.Hour); err != nil || url == "" {
		t.Errorf("unexpected error on GeneratePresignedURL: %v", err)
	}
	if err := svc.DeleteObject(ctx, "new-bucket", "test-key"); err != nil {
		t.Errorf("unexpected error on DeleteObject: %v", err)
	}
	if err := svc.DeleteBucket(ctx, "new-bucket"); err != nil {
		t.Errorf("unexpected error on DeleteBucket: %v", err)
	}
}
