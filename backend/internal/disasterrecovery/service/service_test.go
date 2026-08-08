package service

import (
	"context"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/disasterrecovery/model"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type mockRepo struct {
	plans map[uuid.UUID]*model.RecoveryPlan
	jobs  map[uuid.UUID]*model.RecoveryJob
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		plans: make(map[uuid.UUID]*model.RecoveryPlan),
		jobs:  make(map[uuid.UUID]*model.RecoveryJob),
	}
}

func (m *mockRepo) CreatePlan(ctx context.Context, plan *model.RecoveryPlan) error {
	m.plans[plan.ID] = plan
	return nil
}
func (m *mockRepo) GetPlanByID(ctx context.Context, id uuid.UUID) (*model.RecoveryPlan, error) {
	plan, ok := m.plans[id]
	if !ok {
		return nil, gormNotFound()
	}
	return plan, nil
}
func (m *mockRepo) ListPlansByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.RecoveryPlan, error) {
	var res []model.RecoveryPlan
	for _, plan := range m.plans {
		if plan.WorkspaceID == workspaceID {
			res = append(res, *plan)
		}
	}
	return res, nil
}
func (m *mockRepo) CreateJob(ctx context.Context, job *model.RecoveryJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) UpdateJob(ctx context.Context, job *model.RecoveryJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*model.RecoveryJob, error) {
	job, ok := m.jobs[id]
	if !ok {
		return nil, gormNotFound()
	}
	return job, nil
}
func (m *mockRepo) ListJobsByPlan(ctx context.Context, planID uuid.UUID) ([]model.RecoveryJob, error) {
	var res []model.RecoveryJob
	for _, job := range m.jobs {
		if job.PlanID == planID {
			res = append(res, *job)
		}
	}
	return res, nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string { return "record not found" }
func gormNotFound() error           { return notFoundErr{} }

func TestService_CreatePlanAndTriggerRecovery(t *testing.T) {
	repo := newMockRepo()
	eventsBus := events.NewDispatcher(nil)
	svc := NewService(repo, nil, eventsBus, nil)

	wsID := uuid.New()
	req := CreatePlanRequest{
		WorkspaceID:    wsID,
		Name:           "DR Plan - US East to West",
		Description:    "Automated failover and restore test",
		SourceProvider: "aws-s3-east",
		TargetProvider: "aws-s3-west",
		RPOSeconds:     3600,
		RTOSeconds:     1800,
	}

	plan, err := svc.CreatePlan(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error creating plan: %v", err)
	}

	now := time.Now().UTC()
	triggerReq := TriggerRecoveryRequest{
		PlanID:      plan.ID,
		PointInTime: &now,
	}

	job, err := svc.TriggerRecovery(context.Background(), triggerReq)
	if err != nil {
		t.Fatalf("failed to trigger recovery: %v", err)
	}
	if job.Status != model.JobStatusPending {
		t.Fatalf("expected PENDING status, got %s", job.Status)
	}

	err = svc.MarkJobCompleted(context.Background(), job.ID, 250, 102400)
	if err != nil {
		t.Fatalf("failed to mark job completed: %v", err)
	}

	updatedJob, _ := svc.GetJob(context.Background(), job.ID)
	if updatedJob.Status != model.JobStatusCompleted || updatedJob.ObjectsRestored != 250 {
		t.Fatalf("expected completed job with 250 objects restored, got %+v", updatedJob)
	}
}
