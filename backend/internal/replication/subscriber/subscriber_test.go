package subscriber

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/replication/model"
	"github.com/cloudstorex/backend/internal/replication/service"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type mockService struct {
	scheduledCount int
	lastReq        service.ScheduleRequest
}

func (m *mockService) ScheduleReplication(ctx context.Context, req service.ScheduleRequest) (*model.ObjectReplication, error) {
	m.scheduledCount++
	m.lastReq = req
	return &model.ObjectReplication{ID: uuid.New(), Status: model.StatusPending}, nil
}
func (m *mockService) MarkCompleted(ctx context.Context, id uuid.UUID, checksum string, bytesReplicated int64) error {
	return nil
}
func (m *mockService) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return nil
}
func (m *mockService) GetReplicationsByObject(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error) {
	return nil, nil
}

func TestRegisterReplicationSubscriber(t *testing.T) {
	dispatcher := events.NewDispatcher(nil)
	svc := &mockService{}

	err := RegisterReplicationSubscriber(dispatcher, svc, nil)
	if err != nil {
		t.Fatalf("failed to register subscriber: %v", err)
	}

	wsID := uuid.New()
	objID := uuid.New()
	evt := events.NewEvent(events.EventObjectUploaded, wsID, objID, "", map[string]interface{}{
		"bucket_id":         uuid.NewString(),
		"object_key":        "finance/report-2026.pdf",
		"primary_provider":  "aws-s3",
		"replica_providers": []string{"minio", "gcp-gcs"},
		"mode":              "ASYNC",
	})

	err = dispatcher.Publish(context.Background(), evt)
	if err != nil {
		t.Fatalf("failed to publish event: %v", err)
	}

	if svc.scheduledCount != 2 {
		t.Fatalf("expected 2 replications scheduled for 2 replica providers, got %d", svc.scheduledCount)
	}
	if svc.lastReq.PrimaryProvider != "aws-s3" || svc.lastReq.ReplicaProvider != "gcp-gcs" {
		t.Fatalf("unexpected schedule request: %+v", svc.lastReq)
	}
}
