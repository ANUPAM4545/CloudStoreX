package cluster

import "time"

type ProviderHealthState string

const (
	StateHealthy   ProviderHealthState = "HEALTHY"
	StateDegraded  ProviderHealthState = "DEGRADED"
	StateUnhealthy ProviderHealthState = "UNHEALTHY"
	StateOffline   ProviderHealthState = "OFFLINE"
)

type RecoveryState string

const (
	RecoveryStateIdle                RecoveryState = "IDLE"
	RecoveryStateFailoverInProgress  RecoveryState = "FAILOVER_IN_PROGRESS"
	RecoveryStateRecoveryInProgress  RecoveryState = "RECOVERY_IN_PROGRESS"
)

type NodeStatus string

const (
	NodeOnline   NodeStatus = "ONLINE"
	NodeDraining NodeStatus = "DRAINING"
	NodeOffline  NodeStatus = "OFFLINE"
)

// ProviderHealth represents real-time operational metrics for a storage provider.
type ProviderHealth struct {
	ProviderID    string              `json:"provider_id"`
	State         ProviderHealthState `json:"state"`
	Weight        int                 `json:"weight"`
	LatencyMs     int64               `json:"latency_ms"`
	ErrorRate     float64             `json:"error_rate"`
	LastCheckedAt time.Time           `json:"last_checked_at"`
}

// ClusterState represents the unified operational view of all providers and reliability subsystems.
type ClusterState struct {
	ClusterID      string                    `json:"cluster_id"`
	OverallStatus  string                    `json:"overall_status"`
	ProviderStates map[string]ProviderHealth `json:"provider_states"`
	RegionStates   map[string]string         `json:"region_states"`
	RecoveryState  RecoveryState             `json:"recovery_state"`
	ActiveNodes    int                       `json:"active_nodes"`
	LastUpdated    time.Time                 `json:"last_updated"`
}
