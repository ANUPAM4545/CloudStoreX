package storage

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

// ContextKey defines typed context keys for safe context propagation.
type ContextKey string

const (
	CtxKeyRequestID     ContextKey = "request_id"
	CtxKeyUserID        ContextKey = "user_id"
	CtxKeyWorkspaceID   ContextKey = "workspace_id"
	CtxKeyCorrelationID ContextKey = "correlation_id"
)

// ProviderType represents the type of storage provider.
type ProviderType string

const (
	ProviderTypeDefault ProviderType = "default"
	ProviderTypeS3      ProviderType = "s3"
	ProviderTypeGCS     ProviderType = "gcs"
	ProviderTypeMinIO   ProviderType = "minio"
)

// ProviderInfo represents metadata about a storage provider.
type ProviderInfo struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      ProviderType      `json:"type"`
	Region    string            `json:"region,omitempty"`
	Config    map[string]string `json:"-"` // Don't expose sensitive config in JSON
	CreatedAt time.Time         `json:"created_at"`
}

// Bucket represents a logical container for objects.
type Bucket struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	WorkspaceID uuid.UUID         `json:"workspace_id"`
	ProviderID  string            `json:"provider_id"`
	Region      string            `json:"region,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// ObjectMetadata represents metadata and tags associated with an object.
type ObjectMetadata struct {
	ContentType        string            `json:"content_type"`
	ContentEncoding    string            `json:"content_encoding,omitempty"`
	ContentDisposition string            `json:"content_disposition,omitempty"`
	Custom             map[string]string `json:"custom,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
	VersionID          string            `json:"version_id,omitempty"`
	Encryption         string            `json:"encryption,omitempty"`
}

// Object represents a file or blob stored in a Bucket.
type Object struct {
	Key          string          `json:"key"`
	Bucket       string          `json:"bucket"`
	Size         int64           `json:"size"`
	ETag         string          `json:"etag,omitempty"`
	LastModified time.Time       `json:"last_modified"`
	Metadata     *ObjectMetadata `json:"metadata,omitempty"`
}

// StorageRequest represents a unified request for a storage operation.
type StorageRequest struct {
	Op       string          `json:"op"`
	Bucket   string          `json:"bucket"`
	Key      string          `json:"key,omitempty"`
	Size     int64           `json:"size,omitempty"`
	Reader   io.Reader       `json:"-"`
	Metadata *ObjectMetadata `json:"metadata,omitempty"`
	// Additional targets for operations like Copy/Move
	DestBucket string `json:"dest_bucket,omitempty"`
	DestKey    string `json:"dest_key,omitempty"`
}

// StorageResponse represents a standardized response from any storage provider.
type StorageResponse struct {
	Bucket       string          `json:"bucket"`
	Key          string          `json:"key"`
	Size         int64           `json:"size"`
	ETag         string          `json:"etag,omitempty"`
	VersionID    string          `json:"version_id,omitempty"`
	LastModified time.Time       `json:"last_modified"`
	Metadata     *ObjectMetadata `json:"metadata,omitempty"`
	Body         io.ReadCloser   `json:"-"` // For downloads
}

// MultipartUpload represents an in-progress multipart upload session.
type MultipartUpload struct {
	UploadID   string    `json:"upload_id"`
	Bucket     string    `json:"bucket"`
	Key        string    `json:"key"`
	ProviderID string    `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// UploadPart represents a single part in a multipart upload.
type UploadPart struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
	Size       int64  `json:"size"`
}

// Helper to extract string values from context
func GetContextValue(ctx context.Context, key ContextKey) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}
