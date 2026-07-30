package minio

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/cloudstorex/backend/internal/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient defines the subset of MinIO SDK operations used by the provider, enabling easy mocking.
type MinIOClient interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
}

// Config holds MinIO connection and default bucket settings.
type Config struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	DefaultBucket string
	UseSSL        bool
}

type Provider struct {
	client MinIOClient
	config *Config
	log    *slog.Logger
}

// NewProvider creates and initializes a MinIO storage provider adapter.
// It verifies connectivity and ensures the default bucket exists.
func NewProvider(ctx context.Context, cfg *Config, log *slog.Logger) (*Provider, error) {
	if cfg == nil {
		return nil, errors.New("minio config cannot be nil")
	}
	if log == nil {
		log = slog.Default()
	}

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, storage.NewDomainError("init_minio", "minio", "", "", err)
	}

	p := &Provider{
		client: minioClient,
		config: cfg,
		log:    log,
	}

	if err := p.verifyAndInit(ctx); err != nil {
		return nil, err
	}

	return p, nil
}

// NewProviderWithClient creates a MinIO provider with a custom MinIOClient interface (for testing).
func NewProviderWithClient(ctx context.Context, client MinIOClient, cfg *Config, log *slog.Logger) (*Provider, error) {
	if log == nil {
		log = slog.Default()
	}
	p := &Provider{
		client: client,
		config: cfg,
		log:    log,
	}
	if err := p.verifyAndInit(ctx); err != nil {
		return nil, err
	}
	return p, nil
}

// verifyAndInit verifies MinIO connectivity and ensures the default bucket exists.
func (p *Provider) verifyAndInit(ctx context.Context) error {
	bucket := p.config.DefaultBucket
	if bucket == "" {
		bucket = "cloudstorex-default"
	}

	exists, err := p.client.BucketExists(ctx, bucket)
	if err != nil {
		p.log.ErrorContext(ctx, "MinIO connectivity check failed", slog.String("error", err.Error()))
		return storage.NewDomainError("connectivity_check", "minio", bucket, "", p.translateError("BucketExists", bucket, "", err))
	}

	if !exists {
		err := p.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return storage.NewDomainError("make_default_bucket", "minio", bucket, "", p.translateError("CreateBucket", bucket, "", err))
		}
		p.log.InfoContext(ctx, "Created default MinIO bucket", slog.String("bucket", bucket))
	} else {
		p.log.InfoContext(ctx, "MinIO connectivity verified", slog.String("default_bucket", bucket))
	}

	return nil
}

func (p *Provider) translateError(op, bucket, key string, err error) error {
	if err == nil {
		return nil
	}

	var minioErr minio.ErrorResponse
	if errors.As(err, &minioErr) {
		switch minioErr.Code {
		case "NoSuchKey":
			return storage.NewDomainError(op, "minio", bucket, key, storage.ErrObjectNotFound)
		case "NoSuchBucket":
			return storage.NewDomainError(op, "minio", bucket, key, storage.ErrBucketNotFound)
		case "BucketAlreadyOwnedByYou", "BucketAlreadyExists":
			return storage.NewDomainError(op, "minio", bucket, key, storage.ErrBucketExists)
		}
	}

	// Translate HTTP error codes if present
	if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "The specified key does not exist") {
		return storage.NewDomainError(op, "minio", bucket, key, storage.ErrObjectNotFound)
	}

	switch op {
	case "Upload":
		return storage.NewDomainError(op, "minio", bucket, key, storage.ErrUploadFailed)
	case "Download":
		return storage.NewDomainError(op, "minio", bucket, key, storage.ErrDownloadFailed)
	case "Delete":
		return storage.NewDomainError(op, "minio", bucket, key, storage.ErrDeleteFailed)
	}

	return storage.NewDomainError(op, "minio", bucket, key, err)
}

func (p *Provider) logOperation(ctx context.Context, op, bucket, key string, start time.Time, err error) {
	duration := time.Since(start)
	level := slog.LevelInfo
	result := "success"
	var errStr string

	if err != nil {
		level = slog.LevelError
		result = "error"
		errStr = err.Error()
	}

	p.log.Log(ctx, level, "MinIO Provider Operation",
		slog.String("provider", "minio"),
		slog.String("bucket", bucket),
		slog.String("key", key),
		slog.String("op", op),
		slog.Duration("duration", duration),
		slog.String("result", result),
		slog.String("error", errStr),
	)
}

func (p *Provider) Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *storage.ObjectMetadata) (*storage.StorageResponse, error) {
	start := time.Now()
	contentType := "application/octet-stream"
	userMeta := make(map[string]string)
	if meta != nil {
		if meta.ContentType != "" {
			contentType = meta.ContentType
		}
		if meta.Custom != nil {
			userMeta = meta.Custom
		}
	}

	opts := minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: userMeta,
	}

	info, err := p.client.PutObject(ctx, bucket, key, reader, size, opts)
	if err != nil {
		err = p.translateError("Upload", bucket, key, err)
		p.logOperation(ctx, "Upload", bucket, key, start, err)
		return nil, err
	}

	p.logOperation(ctx, "Upload", bucket, key, start, nil)
	return &storage.StorageResponse{
		Bucket:       info.Bucket,
		Key:          info.Key,
		Size:         info.Size,
		ETag:         info.ETag,
		VersionID:    info.VersionID,
		LastModified: info.LastModified,
		Metadata:     meta,
	}, nil
}

func (p *Provider) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	start := time.Now()
	obj, err := p.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		err = p.translateError("Download", bucket, key, err)
		p.logOperation(ctx, "Download", bucket, key, start, err)
		return nil, err
	}

	// StatObject to check existence immediately and avoid late 404 stream error
	stat, err := obj.Stat()
	if err != nil {
		_ = obj.Close()
		err = p.translateError("Download", bucket, key, err)
		p.logOperation(ctx, "Download", bucket, key, start, err)
		return nil, err
	}

	p.logOperation(ctx, "Download", bucket, key, start, nil)
	_ = stat // validated
	return obj, nil
}

func (p *Provider) Delete(ctx context.Context, bucket, key string) error {
	start := time.Now()
	err := p.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		err = p.translateError("Delete", bucket, key, err)
		p.logOperation(ctx, "Delete", bucket, key, start, err)
		return err
	}
	p.logOperation(ctx, "Delete", bucket, key, start, nil)
	return nil
}

func (p *Provider) Exists(ctx context.Context, bucket, key string) (bool, error) {
	start := time.Now()
	_, err := p.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
	if err != nil {
		translated := p.translateError("Exists", bucket, key, err)
		var domainErr *storage.DomainError
		if errors.As(translated, &domainErr) && errors.Is(domainErr.Err, storage.ErrObjectNotFound) {
			p.logOperation(ctx, "Exists", bucket, key, start, nil)
			return false, nil
		}
		if errors.Is(translated, storage.ErrObjectNotFound) {
			p.logOperation(ctx, "Exists", bucket, key, start, nil)
			return false, nil
		}
		p.logOperation(ctx, "Exists", bucket, key, start, translated)
		return false, translated
	}
	p.logOperation(ctx, "Exists", bucket, key, start, nil)
	return true, nil
}

func (p *Provider) ListObjects(ctx context.Context, bucket, prefix string) ([]*storage.Object, error) {
	start := time.Now()
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	var results []*storage.Object
	for objInfo := range p.client.ListObjects(ctx, bucket, opts) {
		if objInfo.Err != nil {
			err := p.translateError("ListObjects", bucket, prefix, objInfo.Err)
			p.logOperation(ctx, "ListObjects", bucket, prefix, start, err)
			return nil, err
		}
		results = append(results, &storage.Object{
			Key:          objInfo.Key,
			Bucket:       bucket,
			Size:         objInfo.Size,
			ETag:         objInfo.ETag,
			LastModified: objInfo.LastModified,
		})
	}

	p.logOperation(ctx, "ListObjects", bucket, prefix, start, nil)
	return results, nil
}

func (p *Provider) CreateBucket(ctx context.Context, bucket string) error {
	start := time.Now()
	err := p.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		err = p.translateError("CreateBucket", bucket, "", err)
		p.logOperation(ctx, "CreateBucket", bucket, "", start, err)
		return err
	}
	p.logOperation(ctx, "CreateBucket", bucket, "", start, nil)
	return nil
}

func (p *Provider) DeleteBucket(ctx context.Context, bucket string) error {
	start := time.Now()
	// Using unsupported feature or simple delete bucket if SDK client supports it
	// For interface compliance, let's return unsupported or unsupported feature
	p.logOperation(ctx, "DeleteBucket", bucket, "", start, storage.ErrUnsupportedFeature)
	return storage.ErrUnsupportedFeature
}

func (p *Provider) ListBuckets(ctx context.Context) ([]*storage.Bucket, error) {
	start := time.Now()
	buckets, err := p.client.ListBuckets(ctx)
	if err != nil {
		err = p.translateError("ListBuckets", "", "", err)
		p.logOperation(ctx, "ListBuckets", "", "", start, err)
		return nil, err
	}
	var results []*storage.Bucket
	for _, b := range buckets {
		results = append(results, &storage.Bucket{Name: b.Name, CreatedAt: b.CreationDate})
	}
	p.logOperation(ctx, "ListBuckets", "", "", start, nil)
	return results, nil
}

// All advanced operations deferred to later epics as per Scope Reduction
func (p *Provider) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "", storage.ErrUnsupportedFeature
}
func (p *Provider) CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) MoveObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) GetObjectMetadata(ctx context.Context, bucket, key string) (*storage.ObjectMetadata, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) SetObjectMetadata(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) error {
	return storage.ErrUnsupportedFeature
}
func (p *Provider) GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error {
	return storage.ErrUnsupportedFeature
}
func (p *Provider) CreateMultipartUpload(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) (*storage.MultipartUpload, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*storage.UploadPart, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*storage.UploadPart) (*storage.StorageResponse, error) {
	return nil, storage.ErrUnsupportedFeature
}
func (p *Provider) AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error {
	return storage.ErrUnsupportedFeature
}
func (p *Provider) ListObjectVersions(ctx context.Context, bucket, key string) ([]*storage.Object, error) {
	return nil, storage.ErrUnsupportedFeature
}

// Ensure interface compliance at compile time.
var _ storage.StorageProvider = (*Provider)(nil)
