package storage

import (
	"context"
	"io"
	"time"
)

// BucketProvider defines bucket lifecycle and management operations.
type BucketProvider interface {
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error
	ListBuckets(ctx context.Context) ([]*Bucket, error)
}

// ObjectProvider defines object lifecycle, streaming I/O, and advanced object management operations.
type ObjectProvider interface {
	// Basic Object Lifecycle
	Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error)
	Download(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucket, key string) error
	Exists(ctx context.Context, bucket, key string) (bool, error)
	ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error)

	// Presigned URLs
	GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)

	// Extended Object Operations (Copy, Move, Metadata, Tagging)
	CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*StorageResponse, error)
	MoveObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*StorageResponse, error)
	GetObjectMetadata(ctx context.Context, bucket, key string) (*ObjectMetadata, error)
	SetObjectMetadata(ctx context.Context, bucket, key string, meta *ObjectMetadata) error
	GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error)
	SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error

	// Multipart Upload Lifecycle
	CreateMultipartUpload(ctx context.Context, bucket, key string, meta *ObjectMetadata) (*MultipartUpload, error)
	UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*UploadPart, error)
	CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*UploadPart) (*StorageResponse, error)
	AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error

	// Object Versioning
	ListObjectVersions(ctx context.Context, bucket, key string) ([]*Object, error)
}

// StorageProvider combines BucketProvider and ObjectProvider into a unified extensible contract.
type StorageProvider interface {
	BucketProvider
	ObjectProvider
}
