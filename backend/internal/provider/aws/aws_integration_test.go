package aws_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/provider/aws"
	"github.com/cloudstorex/backend/internal/storage"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func TestAWSProvider_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	localstackC, err := localstack.Run(ctx, "localstack/localstack:latest")
	if err != nil {
		t.Fatalf("failed to start localstack container: %v", err)
	}
	defer func() {
		if err := localstackC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate localstack container: %v", err)
		}
	}()

	host, err := localstackC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get localstack host: %v", err)
	}
	port, err := localstackC.MappedPort(ctx, "4566")
	if err != nil {
		t.Fatalf("failed to get localstack port: %v", err)
	}

	endpoint := "http://" + host + ":" + port.Port()

	cfg := &aws.Config{
		Region:          "us-east-1",
		AccessKey:       "test",
		SecretKey:       "test",
		Endpoint:        endpoint,
		DefaultBucket:   "integration-bucket",
		ForcePathStyle:  true,
	}

	provider, err := aws.NewProvider(ctx, cfg, nil)
	if err != nil {
		t.Fatalf("failed to create aws provider: %v", err)
	}

	// Wait for localstack to be fully ready
	time.Sleep(2 * time.Second)

	// 1. Upload Object
	data := []byte("hello aws integration test")
	meta := &storage.ObjectMetadata{
		ContentType: "text/plain",
		Custom: map[string]string{
			"custom-header": "test-value",
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
	if string(downloadedData) != "hello aws integration test" {
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
