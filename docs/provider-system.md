# CloudStoreX Pluggable Provider System

CloudStoreX uses an Interface Segregation model to abstract cloud storage providers.

---

## 1. Provider Interfaces (`internal/provider/provider.go`)

### `BucketProvider`
Manages container-level lifecycle operations:
```go
type BucketProvider interface {
    CreateBucket(ctx context.Context, bucket string) error
    DeleteBucket(ctx context.Context, bucket string) error
    ListBuckets(ctx context.Context) ([]string, error)
    BucketExists(ctx context.Context, bucket string) (bool, error)
}
```

### `ObjectProvider`
Manages object-level streaming I/O and metadata:
```go
type ObjectProvider interface {
    UploadObject(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) (*ObjectMetadata, error)
    DownloadObject(ctx context.Context, bucket, key string) (io.ReadCloser, *ObjectMetadata, error)
    DeleteObject(ctx context.Context, bucket, key string) error
    ListObjects(ctx context.Context, bucket, prefix string) ([]*ObjectMetadata, error)
    ObjectExists(ctx context.Context, bucket, key string) (bool, error)
}
```

### `StorageProvider`
Composes `BucketProvider` and `ObjectProvider` into a unified provider driver:
```go
type StorageProvider interface {
    BucketProvider
    ObjectProvider
    ProviderName() string
    HealthCheck(ctx context.Context) error
}
```

---

## 2. Thread-Safe Registry & Factory

- **`Registry` (`internal/provider/registry.go`)**: Implements thread-safe runtime provider registration using `sync.RWMutex`. Allows dynamic provider discovery and hot-swapping without service restarts.
- **`Factory` (`internal/provider/factory.go`)**: Instantiates provider implementations based on environment or database configuration parameters.
