package worker

import (
	"context"
	"encoding/json"
	"testing"

	jobModel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/cloudstorex/backend/internal/verification/model"
	"github.com/cloudstorex/backend/internal/verification/service"
	"github.com/google/uuid"
)

type dummyVerifService struct {
	completed bool
}

func (d *dummyVerifService) ScheduleVerification(ctx context.Context, req service.ScheduleCheckRequest) (*model.ConsistencyJob, error) {
	return nil, nil
}
func (d *dummyVerifService) RecordDiscrepancy(ctx context.Context, jobID, wsID, bucketID, objID uuid.UUID, objKey string, discType model.DiscrepancyType, details string) error {
	return nil
}
func (d *dummyVerifService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, totalChecked, discrepanciesFound int64) error {
	d.completed = true
	return nil
}
func (d *dummyVerifService) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return nil
}
func (d *dummyVerifService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.ConsistencyJob, error) {
	return nil, nil
}
func (d *dummyVerifService) ListDiscrepancies(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]model.DiscrepancyRecord, error) {
	return nil, nil
}

func TestWorker_Handle(t *testing.T) {
	svc := &dummyVerifService{}
	w := NewWorker(svc, nil, nil)

	payload := ConsistencyCheckPayload{
		JobID:       uuid.NewString(),
		WorkspaceID: uuid.NewString(),
		BucketID:    uuid.NewString(),
		ProviderID:  "aws-s3-east",
	}
	data, _ := json.Marshal(payload)

	job := &jobModel.Job{
		ID:      uuid.NewString(),
		Type:    "CONSISTENCY_CHECK",
		Payload: data,
	}

	err := w.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected worker error: %v", err)
	}
	if !svc.completed {
		t.Fatalf("expected consistency check to mark completed")
	}
}
