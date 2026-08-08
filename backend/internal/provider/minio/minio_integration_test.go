package minio_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/provider/minio"
	"github.com/cloudstorex/backend/internal/storage"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMinIOProvider_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000/tcp").WithStartupTimeout(60 * time.Second),
	}
	minioC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start minio container: %v", err)
	}
	defer func() {
		if err := minioC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate minio container: %v", err)
		}
	}()

	host, err := minioC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := minioC.MappedPort(ctx, "9000")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	endpoint := host + ":" + port.Port()

	cfg := &minio.Config{
		Endpoint:      endpoint,
		AccessKey:     "minioadmin",
		SecretKey:     "minioadmin",
		DefaultBucket: "integration-bucket",
		UseSSL:        false,
	}

	provider, err := minio.NewProvider(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("failed to create minio provider: %v", err)
	}

	// 1. Upload Object
	data := []byte("hello integration test")
	meta := &storage.ObjectMetadata{
		ContentType: "text/plain",
		Custom: map[string]string{
			"x-custom-header": "test-value",
		},
	}
	_, err = provider.Upload(ctx, "integration-bucket", "test.txt", bytes.NewReader(data), int64(len(data)), meta)
	if err != nil {
		t.Fatalf("failed to upload object: %v", err)
	}

	// 2. Download Object
	reader, err := provider.Download(ctx, "integration-bucket", "test.txt")
	if err != nil {
		t.Fatalf("failed to download object: %v", err)
	}
	defer reader.Close()
	downloadedData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read downloaded data: %v", err)
	}
	if string(downloadedData) != "hello integration test" {
		t.Fatalf("downloaded data mismatch: got %q", string(downloadedData))
	}

	// 3. Exists
	exists, err := provider.Exists(ctx, "integration-bucket", "test.txt")
	if err != nil || !exists {
		t.Fatalf("expected object to exist: %v", err)
	}

	// 4. Set/Get Tags
	tags := map[string]string{"env": "test"}
	err = provider.SetObjectTags(ctx, "integration-bucket", "test.txt", tags)
	if err != nil {
		t.Fatalf("failed to set object tags: %v", err)
	}
	fetchedTags, err := provider.GetObjectTags(ctx, "integration-bucket", "test.txt")
	if err != nil {
		t.Fatalf("failed to get object tags: %v", err)
	}
	if fetchedTags["env"] != "test" {
		t.Fatalf("expected tag env=test, got %v", fetchedTags["env"])
	}

	// 5. Delete Object
	err = provider.Delete(ctx, "integration-bucket", "test.txt")
	if err != nil {
		t.Fatalf("failed to delete object: %v", err)
	}
	
	// 6. Exists after delete
	exists, err = provider.Exists(ctx, "integration-bucket", "test.txt")
	if err != nil || exists {
		t.Fatalf("expected object to not exist after delete: %v", err)
	}

	// 7. Error Translation verification
	_, err = provider.Download(ctx, "integration-bucket", "missing.txt")
	if err == nil || !strings.Contains(err.Error(), "object not found") {
		t.Fatalf("expected object not found error, got %v", err)
	}
}
