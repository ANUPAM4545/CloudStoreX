package cluster

import (
	"sort"
	"testing"
)

func TestStateManager_UpdateAndGetProviderHealth(t *testing.T) {
	sm := NewStateManager("test-cluster")

	sm.UpdateProviderHealth("aws-s3", StateHealthy, 100, 15, 0.0)
	sm.UpdateProviderHealth("minio", StateUnhealthy, 50, 200, 0.15)

	awsHealth, ok := sm.GetProviderHealth("aws-s3")
	if !ok || awsHealth.State != StateHealthy {
		t.Fatalf("expected aws-s3 to be HEALTHY, got %v", awsHealth.State)
	}

	healthy := sm.GetHealthyProviders()
	if len(healthy) != 1 || healthy[0] != "aws-s3" {
		t.Fatalf("expected only aws-s3 in healthy providers, got %v", healthy)
	}

	state := sm.GetClusterState()
	if state.OverallStatus != "DEGRADED" {
		t.Fatalf("expected DEGRADED cluster status due to minio unhealthy, got %s", state.OverallStatus)
	}

	sm.SetRecoveryState(RecoveryStateFailoverInProgress)
	state = sm.GetClusterState()
	if state.OverallStatus != "CRITICAL" {
		t.Fatalf("expected CRITICAL cluster status when recovery in progress, got %s", state.OverallStatus)
	}
}

func TestStateManager_AllHealthy(t *testing.T) {
	sm := NewStateManager("prod-cluster")
	sm.UpdateProviderHealth("aws-s3", StateHealthy, 100, 10, 0.0)
	sm.UpdateProviderHealth("minio", StateHealthy, 100, 12, 0.0)

	healthy := sm.GetHealthyProviders()
	sort.Strings(healthy)
	if len(healthy) != 2 || healthy[0] != "aws-s3" || healthy[1] != "minio" {
		t.Fatalf("expected both providers healthy, got %v", healthy)
	}

	if sm.GetClusterState().OverallStatus != "HEALTHY" {
		t.Fatalf("expected HEALTHY overall status")
	}
}
