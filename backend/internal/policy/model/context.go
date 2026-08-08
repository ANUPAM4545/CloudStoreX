package model

import "github.com/google/uuid"

// EvaluationContext encapsulates all metadata needed by the Policy Engine
// to make routing, compliance, replication, and cost-optimization decisions.
type EvaluationContext struct {
	WorkspaceID    uuid.UUID
	OrganizationID uuid.UUID
	Bucket         string
	ObjectKey      string
	Size           int64
	MimeType       string
	StorageClass   string
	Region         string
	ComplianceTags []string
	CostOptimized  bool
	Metadata       map[string]string
	Tags           map[string]string
}
