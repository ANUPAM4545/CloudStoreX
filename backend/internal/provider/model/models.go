package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderStatus string
type ProviderHealth string

const (
	StatusPending    ProviderStatus = "PENDING"
	StatusValidating ProviderStatus = "VALIDATING"
	StatusReady      ProviderStatus = "READY"
	StatusDisabled   ProviderStatus = "DISABLED"
	StatusFailed     ProviderStatus = "FAILED"

	HealthHealthy     ProviderHealth = "HEALTHY"
	HealthDegraded    ProviderHealth = "DEGRADED"
	HealthUnavailable ProviderHealth = "UNAVAILABLE"
	HealthUnknown     ProviderHealth = "UNKNOWN"
)

// Provider represents a storage provider registered in the system.
type Provider struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey"`
	WorkspaceID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_workspace_provider_name"`
	ProviderName     string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_workspace_provider_name"`
	ProviderType     string         `gorm:"type:varchar(100);not null"` // e.g., MINIO, AWS_S3
	Endpoint         string         `gorm:"type:varchar(512)"`          // Optional
	Region           string         `gorm:"type:varchar(100)"`
	BucketPrefix     string         `gorm:"type:varchar(100)"`
	IsDefault        bool           `gorm:"default:false;index"`
	IsEnabled        bool           `gorm:"default:true;index"`
	Status           ProviderStatus `gorm:"type:varchar(50);default:'PENDING'"`
	Health           ProviderHealth `gorm:"type:varchar(50);default:'UNKNOWN'"`
	LastHealthCheck  *time.Time     
	LatencyMs        int64          `gorm:"default:0"`
	LastError        string         `gorm:"type:text"`
	CredentialRef    string         `gorm:"type:varchar(255)"` // Reference to secrets manager
	CapabilitiesJSON string         `gorm:"type:text"`         // JSON serialized ProviderCapabilities
	CreatedAt        time.Time      `gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// ProviderCapabilities describes features supported by a given provider implementation.
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
