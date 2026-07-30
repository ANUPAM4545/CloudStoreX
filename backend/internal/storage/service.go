package storage

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/shared/logger"
)

// Service defines the application-level storage contract.
// No HTTP handlers should communicate directly with providers; they must use this Service.
type Service interface {
	UploadObject(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error)
	DownloadObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucket, key string) error
	ObjectExists(ctx context.Context, bucket, key string) (bool, error)
	ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error)
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error
	ListBuckets(ctx context.Context) ([]*Bucket, error)
	GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)
}

type defaultService struct {
	router Router
	log    *slog.Logger
}

// NewService creates a new high-level Storage Service.
func NewService(router Router) Service {
	l := logger.Log
	if l == nil {
		l = slog.Default()
	}
	return &defaultService{
		router: router,
		log:    l,
	}
}

func (s *defaultService) logOperation(ctx context.Context, op, bucket, key string, start time.Time, err error) {
	if s.log == nil {
		s.log = slog.Default()
	}
	duration := time.Since(start)
	reqID := GetContextValue(ctx, CtxKeyRequestID)
	workspaceID := GetContextValue(ctx, CtxKeyWorkspaceID)

	level := slog.LevelInfo
	result := "success"
	var errStr string

	if err != nil {
		level = slog.LevelError
		result = "error"
		errStr = err.Error()
	}

	s.log.Log(ctx, level, "Storage Operation",
		slog.String("op", op),
		slog.String("bucket", bucket),
		slog.String("key", key),
		slog.String("request_id", reqID),
		slog.String("workspace_id", workspaceID),
		slog.Duration("duration", duration),
		slog.String("result", result),
		slog.String("error", errStr),
	)
}

func (s *defaultService) UploadObject(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error) {
	start := time.Now()
	res, err := s.router.Upload(ctx, bucket, key, reader, size, meta)
	s.logOperation(ctx, "UploadObject", bucket, key, start, err)
	return res, err
}

func (s *defaultService) DownloadObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	start := time.Now()
	res, err := s.router.Download(ctx, bucket, key)
	s.logOperation(ctx, "DownloadObject", bucket, key, start, err)
	return res, err
}

func (s *defaultService) DeleteObject(ctx context.Context, bucket, key string) error {
	start := time.Now()
	err := s.router.Delete(ctx, bucket, key)
	s.logOperation(ctx, "DeleteObject", bucket, key, start, err)
	return err
}

func (s *defaultService) ObjectExists(ctx context.Context, bucket, key string) (bool, error) {
	start := time.Now()
	exists, err := s.router.Exists(ctx, bucket, key)
	s.logOperation(ctx, "ObjectExists", bucket, key, start, err)
	return exists, err
}

func (s *defaultService) ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error) {
	start := time.Now()
	objs, err := s.router.ListObjects(ctx, bucket, prefix)
	s.logOperation(ctx, "ListObjects", bucket, prefix, start, err)
	return objs, err
}

func (s *defaultService) CreateBucket(ctx context.Context, bucket string) error {
	start := time.Now()
	err := s.router.CreateBucket(ctx, bucket)
	s.logOperation(ctx, "CreateBucket", bucket, "", start, err)
	return err
}

func (s *defaultService) DeleteBucket(ctx context.Context, bucket string) error {
	start := time.Now()
	err := s.router.DeleteBucket(ctx, bucket)
	s.logOperation(ctx, "DeleteBucket", bucket, "", start, err)
	return err
}

func (s *defaultService) ListBuckets(ctx context.Context) ([]*Bucket, error) {
	start := time.Now()
	buckets, err := s.router.ListBuckets(ctx)
	s.logOperation(ctx, "ListBuckets", "", "", start, err)
	return buckets, err
}

func (s *defaultService) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	start := time.Now()
	url, err := s.router.GeneratePresignedURL(ctx, bucket, key, expiration)
	s.logOperation(ctx, "GeneratePresignedURL", bucket, key, start, err)
	return url, err
}
