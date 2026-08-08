package worker

import (
	"context"
	"encoding/json"
	"testing"

	jobModel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/cloudstorex/backend/internal/replication/model"
	"github.com/cloudstorex/backend/internal/replication/service"
	"github.com/google/uuid"
)

type dummyService struct {
	completed bool
}

func (d *dummyService) ScheduleReplication(ctx context.Context, req service.ScheduleRequest) (*model.ObjectReplication, error) {
	return nil, nil
}
func (d *dummyService) MarkCompleted(ctx context.Context, id uuid.UUID, checksum string, bytesReplicated int64) error {
	d.completed = true
	return nil
}
func (d *dummyService) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return nil
}
func (d *dummyService) GetReplicationsByObject(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error) {
	return nil, nil
}

func TestReplicationWorker_Handle(t *testing.T) {
	svc := &dummyService{}
	worker := NewWorker(svc, nil, nil)

	payload := ReplicationJobPayload{
		ReplicationID:   uuid.NewString(),
		WorkspaceID:     uuid.NewString(),
		BucketID:        uuid.NewString(),
		ObjectID:        uuid.NewString(),
		ObjectKey:       "enterprise-report.pdf",
		PrimaryProvider: "aws-s3",
		ReplicaProvider: "minio",
	}
	data, _ := json.Marshal(payload)

	job := &jobModel.Job{
		ID:      uuid.NewString(),
		Type:    "REPLICATE_OBJECT",
		Payload: data,
	}

	err := worker.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected worker handle error: %v", err)
	}
	if !svc.completed {
		t.Fatalf("expected replication to be marked completed")
	}
}
