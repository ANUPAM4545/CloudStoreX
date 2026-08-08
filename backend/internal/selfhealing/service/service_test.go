package service

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/cloudstorex/backend/internal/selfhealing/model"
	"github.com/google/uuid"
)

type mockRepo struct {
	jobs map[uuid.UUID]*model.RepairJob
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		jobs: make(map[uuid.UUID]*model.RepairJob),
	}
}

func (m *mockRepo) CreateJob(ctx context.Context, job *model.RepairJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) UpdateJob(ctx context.Context, job *model.RepairJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*model.RepairJob, error) {
	job, ok := m.jobs[id]
	if !ok {
		return nil, gormNotFound()
	}
	return job, nil
}
func (m *mockRepo) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error) {
	var res []model.RepairJob
	for _, j := range m.jobs {
		if j.WorkspaceID == workspaceID {
			res = append(res, *j)
		}
	}
	return res, nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string { return "record not found" }
func gormNotFound() error           { return notFoundErr{} }

func TestService_ScheduleAndCompleteRepair(t *testing.T) {
	repo := newMockRepo()
	eventsBus := events.NewDispatcher(nil)
	svc := NewService(repo, nil, eventsBus, nil)

	wsID := uuid.New()
	bucketID := uuid.New()
	objID := uuid.New()

	req := ScheduleRepairRequest{
		WorkspaceID:    wsID,
		BucketID:       bucketID,
		ObjectID:       objID,
		ObjectKey:      "archive/doc.pdf",
		TargetProvider: "aws-s3-east",
		SourceProvider: "aws-s3-west",
		SourceType:     model.SourceReplica,
	}

	job, err := svc.ScheduleRepair(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected schedule error: %v", err)
	}

	err = svc.MarkJobCompleted(context.Background(), job.ID, 102400)
	if err != nil {
		t.Fatalf("unexpected completion error: %v", err)
	}

	updated, _ := svc.ListJobs(context.Background(), wsID)
	if len(updated) != 1 || updated[0].Status != model.StatusCompleted || updated[0].BytesRepaired != 102400 {
		t.Fatalf("expected completed repair job with 102400 bytes, got %+v", updated)
	}
}
