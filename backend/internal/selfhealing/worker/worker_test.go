package worker

import (
	"context"
	"encoding/json"
	"testing"

	jobModel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/cloudstorex/backend/internal/selfhealing/model"
	"github.com/cloudstorex/backend/internal/selfhealing/service"
	"github.com/google/uuid"
)

type dummyRepairService struct {
	completed bool
}

func (d *dummyRepairService) ScheduleRepair(ctx context.Context, req service.ScheduleRepairRequest) (*model.RepairJob, error) {
	return nil, nil
}
func (d *dummyRepairService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, bytesRepaired int64) error {
	d.completed = true
	return nil
}
func (d *dummyRepairService) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return nil
}
func (d *dummyRepairService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error) {
	return nil, nil
}

func TestWorker_Handle(t *testing.T) {
	svc := &dummyRepairService{}
	w := NewWorker(svc, nil, nil)

	payload := SelfHealPayload{
		RepairJobID:    uuid.NewString(),
		WorkspaceID:    uuid.NewString(),
		BucketID:       uuid.NewString(),
		ObjectID:       uuid.NewString(),
		ObjectKey:      "archive/doc.pdf",
		TargetProvider: "aws-s3-east",
		SourceProvider: "aws-s3-west",
	}
	data, _ := json.Marshal(payload)

	job := &jobModel.Job{
		ID:      uuid.NewString(),
		Type:    "SELF_HEAL_JOB",
		Payload: data,
	}

	err := w.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected worker error: %v", err)
	}
	if !svc.completed {
		t.Fatalf("expected self healing job to mark completed")
	}
}
