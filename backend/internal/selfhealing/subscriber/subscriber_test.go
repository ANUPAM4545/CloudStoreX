package subscriber

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/cloudstorex/backend/internal/selfhealing/model"
	"github.com/cloudstorex/backend/internal/selfhealing/service"
	"github.com/google/uuid"
)

type dummyRepairService struct {
	scheduledCount int
}

func (d *dummyRepairService) ScheduleRepair(ctx context.Context, req service.ScheduleRepairRequest) (*model.RepairJob, error) {
	d.scheduledCount++
	return &model.RepairJob{ID: uuid.New(), Status: model.StatusPending}, nil
}
func (d *dummyRepairService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, bytesRepaired int64) error {
	return nil
}
func (d *dummyRepairService) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return nil
}
func (d *dummyRepairService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error) {
	return nil, nil
}

func TestRegisterSelfHealingSubscriber(t *testing.T) {
	dispatcher := events.NewDispatcher(nil)
	svc := &dummyRepairService{}

	err := RegisterSelfHealingSubscriber(dispatcher, svc, nil)
	if err != nil {
		t.Fatalf("failed to register subscriber: %v", err)
	}

	wsID := uuid.New()
	objID := uuid.New()
	evt := events.NewEvent(events.EventConsistencyCheckFailed, wsID, objID, "aws-s3-east", map[string]interface{}{
		"bucket_id":        uuid.NewString(),
		"object_key":       "archive/doc.pdf",
		"discrepancy_type": "CHECKSUM_MISMATCH",
		"details":          "hash mismatch",
	})

	err = dispatcher.Publish(context.Background(), evt)
	if err != nil {
		t.Fatalf("failed to publish event: %v", err)
	}

	if svc.scheduledCount != 1 {
		t.Fatalf("expected 1 repair scheduled from consistency event, got %d", svc.scheduledCount)
	}
}
