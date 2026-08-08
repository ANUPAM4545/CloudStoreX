package service

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/replication/model"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type mockRepo struct {
	replications map[uuid.UUID]*model.ObjectReplication
}

func newMockRepo() *mockRepo {
	return &mockRepo{replications: make(map[uuid.UUID]*model.ObjectReplication)}
}

func (m *mockRepo) Create(ctx context.Context, repl *model.ObjectReplication) error {
	m.replications[repl.ID] = repl
	return nil
}

func (m *mockRepo) Update(ctx context.Context, repl *model.ObjectReplication) error {
	m.replications[repl.ID] = repl
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.ObjectReplication, error) {
	repl, ok := m.replications[id]
	if !ok {
		return nil, gormNotFound()
	}
	return repl, nil
}

func (m *mockRepo) ListByObjectID(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error) {
	var res []model.ObjectReplication
	for _, repl := range m.replications {
		if repl.ObjectID == objectID {
			res = append(res, *repl)
		}
	}
	return res, nil
}

func (m *mockRepo) ListByStatus(ctx context.Context, status model.ReplicationStatus, limit int) ([]model.ObjectReplication, error) {
	var res []model.ObjectReplication
	for _, repl := range m.replications {
		if repl.Status == status {
			res = append(res, *repl)
		}
	}
	return res, nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string { return "record not found" }
func gormNotFound() error           { return notFoundErr{} }

func TestService_ScheduleReplication(t *testing.T) {
	repo := newMockRepo()
	eventsBus := events.NewDispatcher(nil)

	var eventCount int
	eventsBus.Subscribe(events.EventReplicationStarted, func(ctx context.Context, event events.Event) error {
		eventCount++
		return nil
	})

	svc := NewService(repo, nil, eventsBus, nil)

	req := ScheduleRequest{
		WorkspaceID:     uuid.New(),
		BucketID:        uuid.New(),
		ObjectID:        uuid.New(),
		ObjectKey:       "enterprise/data.pdf",
		PrimaryProvider: "aws-s3",
		ReplicaProvider: "minio",
		Mode:            model.ModeAsync,
	}

	repl, err := svc.ScheduleReplication(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error scheduling replication: %v", err)
	}

	if repl.Status != model.StatusPending {
		t.Fatalf("expected status PENDING, got %v", repl.Status)
	}

	if eventCount != 1 {
		t.Fatalf("expected 1 ReplicationStarted event, got %d", eventCount)
	}
}

func TestService_MarkCompletedAndFailed(t *testing.T) {
	repo := newMockRepo()
	eventsBus := events.NewDispatcher(nil)
	svc := NewService(repo, nil, eventsBus, nil)

	req := ScheduleRequest{
		WorkspaceID:     uuid.New(),
		BucketID:        uuid.New(),
		ObjectID:        uuid.New(),
		ObjectKey:       "test-obj",
		PrimaryProvider: "aws-s3",
		ReplicaProvider: "minio",
	}

	repl, _ := svc.ScheduleReplication(context.Background(), req)

	err := svc.MarkCompleted(context.Background(), repl.ID, "sha256-hash", 4096)
	if err != nil {
		t.Fatalf("failed to mark completed: %v", err)
	}

	updated, _ := repo.GetByID(context.Background(), repl.ID)
	if updated.Status != model.StatusCompleted || updated.Checksum != "sha256-hash" {
		t.Fatalf("expected COMPLETED with checksum, got %v / %v", updated.Status, updated.Checksum)
	}

	_ = svc.MarkFailed(context.Background(), repl.ID, "replica bucket full")
	updated, _ = repo.GetByID(context.Background(), repl.ID)
	if updated.Status != model.StatusFailed || updated.ErrorMessage != "replica bucket full" {
		t.Fatalf("expected FAILED with error message")
	}
}
