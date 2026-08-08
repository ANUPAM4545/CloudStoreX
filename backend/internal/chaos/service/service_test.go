package service

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/chaos/model"
	"github.com/cloudstorex/backend/internal/cluster"
	"github.com/google/uuid"
)

type mockRepo struct {
	exps map[uuid.UUID]*model.ChaosExperiment
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		exps: make(map[uuid.UUID]*model.ChaosExperiment),
	}
}

func (m *mockRepo) CreateExperiment(ctx context.Context, exp *model.ChaosExperiment) error {
	m.exps[exp.ID] = exp
	return nil
}
func (m *mockRepo) UpdateExperiment(ctx context.Context, exp *model.ChaosExperiment) error {
	m.exps[exp.ID] = exp
	return nil
}
func (m *mockRepo) GetExperimentByID(ctx context.Context, id uuid.UUID) (*model.ChaosExperiment, error) {
	exp, ok := m.exps[id]
	if !ok {
		return nil, gormNotFound()
	}
	return exp, nil
}
func (m *mockRepo) ListExperimentsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.ChaosExperiment, error) {
	var res []model.ChaosExperiment
	for _, e := range m.exps {
		if e.WorkspaceID == workspaceID {
			res = append(res, *e)
		}
	}
	return res, nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string { return "record not found" }
func gormNotFound() error           { return notFoundErr{} }

func TestService_RunAndCompleteExperiment(t *testing.T) {
	repo := newMockRepo()
	stateMgr := cluster.NewStateManager("test-cluster")
	svc := NewService(repo, stateMgr, nil)

	wsID := uuid.New()
	req := RunExperimentRequest{
		WorkspaceID:     wsID,
		Name:            "Simulate AWS S3 East Outage",
		ExperimentType:  model.TypeProviderOutage,
		TargetProvider:  "aws-s3-east",
		DurationSeconds: 120,
		TargetRPO:       300,
		TargetRTO:       60,
	}

	exp, err := svc.RunExperiment(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if exp.Status != model.StatusRunning {
		t.Fatalf("expected RUNNING status, got %s", exp.Status)
	}

	// Verify provider was set UNHEALTHY in state manager
	health, ok := stateMgr.GetProviderHealth("aws-s3-east")
	if !ok || health.State != cluster.StateUnhealthy {
		t.Fatalf("expected aws-s3-east to be UNHEALTHY during outage experiment, got %+v", health)
	}

	// Complete experiment within target SLA (RPO=45s <= 300s, RTO=30s <= 60s)
	err = svc.CompleteExperiment(context.Background(), exp.ID, 45, 30)
	if err != nil {
		t.Fatalf("unexpected completion error: %v", err)
	}

	updated, _ := svc.GetExperiment(context.Background(), exp.ID)
	if !updated.Passed || updated.Status != model.StatusCompleted {
		t.Fatalf("expected PASSED completed experiment, got %+v", updated)
	}

	// Verify provider health restored
	health, ok = stateMgr.GetProviderHealth("aws-s3-east")
	if !ok || health.State != cluster.StateHealthy {
		t.Fatalf("expected aws-s3-east to be restored HEALTHY after experiment, got %+v", health)
	}
}
