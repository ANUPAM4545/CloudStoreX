package service

import (
	"context"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/cluster"
	"github.com/cloudstorex/backend/internal/failover/model"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

func TestService_FailoverAndFailback(t *testing.T) {
	stateMgr := cluster.NewStateManager("test-cluster")
	eventsBus := events.NewDispatcher(nil)

	// Set primary unhealthy and candidates healthy
	stateMgr.UpdateProviderHealth("aws-s3-east", cluster.StateUnhealthy, 100, 5000, 0.5)
	stateMgr.UpdateProviderHealth("aws-s3-west", cluster.StateHealthy, 90, 120, 0.0)

	svc := NewService(stateMgr, eventsBus, FailoverConfig{CooldownPeriod: 10 * time.Millisecond}, nil)
	wsID := uuid.New()

	evt, err := svc.CheckAndFailover(context.Background(), wsID, "aws-s3-east", []string{"aws-s3-west"})
	if err != nil {
		t.Fatalf("unexpected failover error: %v", err)
	}

	if evt.NewPrimaryProvider != "aws-s3-west" || evt.Status != model.StatusActive {
		t.Fatalf("unexpected failover event: %+v", evt)
	}

	active := svc.GetActiveFailovers()
	if len(active) != 1 {
		t.Fatalf("expected 1 active failover, got %d", len(active))
	}

	// Try failback before primary is healthy -> should fail
	_, err = svc.CheckAndFailback(context.Background(), evt.ID)
	if err == nil {
		t.Fatalf("expected failback error when primary is still unhealthy")
	}

	// Recover primary provider health
	stateMgr.UpdateProviderHealth("aws-s3-east", cluster.StateHealthy, 100, 80, 0.0)

	fbEvt, err := svc.CheckAndFailback(context.Background(), evt.ID)
	if err != nil {
		t.Fatalf("unexpected failback error: %v", err)
	}
	if fbEvt.Status != model.StatusFailback {
		t.Fatalf("expected status FAILBACK_COMPLETED, got %s", fbEvt.Status)
	}

	active = svc.GetActiveFailovers()
	if len(active) != 0 {
		t.Fatalf("expected 0 active failovers after failback, got %d", len(active))
	}
}
