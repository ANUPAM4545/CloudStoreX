package storage

import (
	"bytes"
	"context"
	"testing"
)

// BenchmarkStorage_UploadObject benchmarks memory allocations and throughput for object upload.
func BenchmarkStorage_UploadObject(b *testing.B) {
	mockProv := &mockProvider{}
	router := &mockRouter{StorageProvider: mockProv}
	svc := NewService(router, &mockMetadataService{}, nil)

	ctx := context.WithValue(context.Background(), CtxKeyWorkspaceID, "ws-67890")
	data := []byte("hello cloudstorex performance benchmark")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		_, _ = svc.UploadObject(ctx, "bench-bucket", "bench-key", reader, int64(len(data)), &ObjectMetadata{})
	}
}

// BenchmarkStorage_UploadObjectParallel benchmarks multi-goroutine upload scaling under concurrent load.
func BenchmarkStorage_UploadObjectParallel(b *testing.B) {
	mockProv := &mockProvider{}
	router := &mockRouter{StorageProvider: mockProv}
	svc := NewService(router, &mockMetadataService{}, nil)

	ctx := context.WithValue(context.Background(), CtxKeyWorkspaceID, "ws-67890")
	data := []byte("hello cloudstorex parallel benchmark")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reader := bytes.NewReader(data)
			_, _ = svc.UploadObject(ctx, "bench-bucket", "bench-key-parallel", reader, int64(len(data)), &ObjectMetadata{})
		}
	})
}
