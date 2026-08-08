package cluster

import (
	"sync"
	"time"
)

// StateManager defines the central authority for real-time cluster and provider health state.
type StateManager interface {
	UpdateProviderHealth(providerID string, state ProviderHealthState, weight int, latencyMs int64, errorRate float64)
	GetProviderHealth(providerID string) (ProviderHealth, bool)
	GetHealthyProviders() []string
	SetRecoveryState(state RecoveryState)
	GetClusterState() ClusterState
}

type defaultStateManager struct {
	mu             sync.RWMutex
	clusterID      string
	providerStates map[string]ProviderHealth
	regionStates   map[string]string
	recoveryState  RecoveryState
	lastUpdated    time.Time
}

// NewStateManager creates a new concurrency-safe cluster state manager.
func NewStateManager(clusterID string) StateManager {
	if clusterID == "" {
		clusterID = "cloudstorex-prod-cluster"
	}
	return &defaultStateManager{
		clusterID:      clusterID,
		providerStates: make(map[string]ProviderHealth),
		regionStates:   make(map[string]string),
		recoveryState:  RecoveryStateIdle,
		lastUpdated:    time.Now().UTC(),
	}
}

func (s *defaultStateManager) UpdateProviderHealth(providerID string, state ProviderHealthState, weight int, latencyMs int64, errorRate float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.providerStates[providerID] = ProviderHealth{
		ProviderID:    providerID,
		State:         state,
		Weight:        weight,
		LatencyMs:     latencyMs,
		ErrorRate:     errorRate,
		LastCheckedAt: time.Now().UTC(),
	}
	s.lastUpdated = time.Now().UTC()
}

func (s *defaultStateManager) GetProviderHealth(providerID string) (ProviderHealth, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	h, ok := s.providerStates[providerID]
	return h, ok
}

func (s *defaultStateManager) GetHealthyProviders() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var healthy []string
	for id, h := range s.providerStates {
		if h.State == StateHealthy {
			healthy = append(healthy, id)
		}
	}
	return healthy
}

func (s *defaultStateManager) SetRecoveryState(state RecoveryState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recoveryState = state
	s.lastUpdated = time.Now().UTC()
}

func (s *defaultStateManager) GetClusterState() ClusterState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	overall := "HEALTHY"
	for _, h := range s.providerStates {
		if h.State == StateUnhealthy || h.State == StateOffline {
			overall = "DEGRADED"
			break
		}
	}
	if s.recoveryState != RecoveryStateIdle {
		overall = "CRITICAL"
	}

	copyProviders := make(map[string]ProviderHealth, len(s.providerStates))
	for k, v := range s.providerStates {
		copyProviders[k] = v
	}

	return ClusterState{
		ClusterID:      s.clusterID,
		OverallStatus:  overall,
		ProviderStates: copyProviders,
		RegionStates:   s.regionStates,
		RecoveryState:  s.recoveryState,
		ActiveNodes:    1,
		LastUpdated:    s.lastUpdated,
	}
}
