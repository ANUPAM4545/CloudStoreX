package worker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	jobModel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/cloudstorex/backend/internal/disasterrecovery/model"
	"github.com/cloudstorex/backend/internal/disasterrecovery/service"
	"github.com/google/uuid"
)

type dummyDRService struct {
	completed bool
}

func (d *dummyDRService) CreatePlan(ctx context.Context, req service.CreatePlanRequest) (*model.RecoveryPlan, error) {
	return nil, nil
}
func (d *dummyDRService) GetPlan(ctx context.Context, id uuid.UUID) (*model.RecoveryPlan, error) {
	return nil, nil
}
func (d *dummyDRService) ListPlans(ctx context.Context, workspaceID uuid.UUID) ([]model.RecoveryPlan, error) {
	return nil, nil
}
func (d *dummyDRService) TriggerRecovery(ctx context.Context, req service.TriggerRecoveryRequest) (*model.RecoveryJob, error) {
	return nil, nil
}
func (d *dummyDRService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, objectsRestored, bytesRestored int64) error {
	d.completed = true
	return nil
}
func (d *dummyDRService) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return nil
}
func (d *dummyDRService) GetJob(ctx context.Context, id uuid.UUID) (*model.RecoveryJob, error) {
	return nil, nil
}

func TestRestoreWorker_Handle(t *testing.T) {
	svc := &dummyDRService{}
	worker := NewWorker(svc, nil, nil)

	now := time.Now().UTC()
	payload := RestoreJobPayload{
		JobID:          uuid.NewString(),
		PlanID:         uuid.NewString(),
		WorkspaceID:    uuid.NewString(),
		SourceProvider: "aws-s3-east",
		TargetProvider: "aws-s3-west",
		PointInTime:    &now,
	}
	data, _ := json.Marshal(payload)

	job := &jobModel.Job{
		ID:      uuid.NewString(),
		Type:    "RESTORE_JOB",
		Payload: data,
	}

	err := worker.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected restore worker error: %v", err)
	}
	if !svc.completed {
		t.Fatalf("expected restore job to be marked completed")
	}
}
