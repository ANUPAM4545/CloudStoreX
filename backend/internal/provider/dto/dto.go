package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProviderDTO struct {
	ID               uuid.UUID           `json:"id"`
	WorkspaceID      uuid.UUID           `json:"workspace_id"`
	ProviderName     string              `json:"provider_name"`
	ProviderType     string              `json:"provider_type"`
	Endpoint         string              `json:"endpoint,omitempty"`
	Region           string              `json:"region,omitempty"`
	BucketPrefix     string              `json:"bucket_prefix,omitempty"`
	IsDefault        bool                `json:"is_default"`
	IsEnabled        bool                `json:"is_enabled"`
	Status           string              `json:"status"`
	Health           string              `json:"health"`
	LastHealthCheck  *time.Time          `json:"last_health_check,omitempty"`
	LatencyMs        int64               `json:"latency_ms"`
	LastError        string              `json:"last_error,omitempty"`
	Capabilities     ProviderCapabilities `json:"capabilities"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}

type ProviderCapabilities struct {
	MultipartUpload bool `json:"multipart_upload"`
	Versioning      bool `json:"versioning"`
	ObjectTagging   bool `json:"object_tagging"`
	PreSignedURLs   bool `json:"pre_signed_urls"`
	Encryption      bool `json:"encryption"`
	ObjectLock      bool `json:"object_lock"`
	LifecycleRules  bool `json:"lifecycle_rules"`
	Replication     bool `json:"replication"`
}

type CreateProviderRequest struct {
	WorkspaceID      uuid.UUID            `json:"workspace_id" binding:"required"`
	ProviderName     string               `json:"provider_name" binding:"required"`
	ProviderType     string               `json:"provider_type" binding:"required"`
	Endpoint         string               `json:"endpoint"`
	Region           string               `json:"region"`
	BucketPrefix     string               `json:"bucket_prefix"`
	IsDefault        bool                 `json:"is_default"`
	CredentialRef    string               `json:"credential_ref"`
	Capabilities     ProviderCapabilities `json:"capabilities"`
}

type UpdateProviderRequest struct {
	Endpoint         *string               `json:"endpoint"`
	Region           *string               `json:"region"`
	BucketPrefix     *string               `json:"bucket_prefix"`
	CredentialRef    *string               `json:"credential_ref"`
	Capabilities     *ProviderCapabilities `json:"capabilities"`
}
