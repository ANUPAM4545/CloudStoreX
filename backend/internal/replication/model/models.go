package model

import (
	"time"

	"github.com/google/uuid"
)

type ReplicationStatus string

const (
	StatusPending     ReplicationStatus = "PENDING"
	StatusReplicating ReplicationStatus = "REPLICATING"
	StatusCompleted   ReplicationStatus = "COMPLETED"
	StatusFailed      ReplicationStatus = "FAILED"
)

type ReplicationMode string

const (
	ModeAsync  ReplicationMode = "ASYNC"
	ModeSync   ReplicationMode = "SYNC"
	ModeOneWay ReplicationMode = "ONE_WAY"
)

// ObjectReplication tracks the replication state of an object across primary and secondary providers/regions.
type ObjectReplication struct {
	ID              uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkspaceID     uuid.UUID         `gorm:"type:uuid;not null;index" json:"workspace_id"`
	BucketID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"bucket_id"`
	ObjectID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"object_id"`
	ObjectKey       string            `gorm:"type:varchar(1024);not null" json:"object_key"`
	PrimaryProvider string            `gorm:"type:varchar(100);not null" json:"primary_provider"`
	ReplicaProvider string            `gorm:"type:varchar(100);not null;index" json:"replica_provider"`
	Mode            ReplicationMode   `gorm:"type:varchar(20);not null;default:'ASYNC'" json:"mode"`
	Status          ReplicationStatus `gorm:"type:varchar(20);not null;default:'PENDING';index" json:"status"`
	Checksum        string            `gorm:"type:varchar(64)" json:"checksum"`
	BytesReplicated int64             `gorm:"not null;default:0" json:"bytes_replicated"`
	ReplicatedAt    *time.Time        `json:"replicated_at"`
	LastVerifiedAt  *time.Time        `json:"last_verified_at"`
	ErrorMessage    string            `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}
