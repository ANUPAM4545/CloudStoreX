package dto

import "time"

// SearchQuery represents the search parameters for the Metadata Service.
type SearchQuery struct {
	WorkspaceID  string            `json:"workspace_id,omitempty"`
	BucketID     string            `json:"bucket_id,omitempty"`
	Prefix       string            `json:"prefix,omitempty"`
	MimeType     string            `json:"mime_type,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	OwnerID      string            `json:"owner_id,omitempty"`
	ProviderID   string            `json:"provider_id,omitempty"`
	StorageClass string            `json:"storage_class,omitempty"`
	Status       string            `json:"status,omitempty"`
	Limit        int               `json:"limit,omitempty"`
	Offset       int               `json:"offset,omitempty"`
}

// SearchResponse represents paginated metadata search results.
type SearchResponse struct {
	Total   int64           `json:"total"`
	Results []ObjectMetaDTO `json:"results"`
}

// ObjectMetaDTO represents the client-facing view of a metadata object.
type ObjectMetaDTO struct {
	ID                string            `json:"id"`
	BucketID          string            `json:"bucket_id"`
	ObjectKey         string            `json:"object_key"`
	ProviderObjectKey string            `json:"provider_object_key"`
	SizeBytes         int64             `json:"size_bytes"`
	MimeType          string            `json:"mime_type"`
	ETag              string            `json:"etag"`
	ProviderID        string            `json:"provider_id"`
	VersionID         string            `json:"version_id,omitempty"`
	OwnerID           string            `json:"owner_id,omitempty"`
	StorageClass      string            `json:"storage_class"`
	Status            string            `json:"status"`
	Tags              map[string]string `json:"tags,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// ObjectVersionDTO represents a specific version of an object.
type ObjectVersionDTO struct {
	ID                string    `json:"id"`
	ObjectID          string    `json:"object_id"`
	VersionNumber     int       `json:"version_number"`
	ProviderID        string    `json:"provider_id"`
	ProviderObjectKey string    `json:"provider_object_key"`
	SizeBytes         int64     `json:"size_bytes"`
	ETag              string    `json:"etag"`
	IsCurrent         bool      `json:"is_current"`
	CreatedAt         time.Time `json:"created_at"`
}
