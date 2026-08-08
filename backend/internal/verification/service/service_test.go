package service

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/cloudstorex/backend/internal/verification/model"
	"github.com/google/uuid"
)

type mockRepo struct {
	jobs  map[uuid.UUID]*model.ConsistencyJob
	discs map[uuid.UUID]*model.DiscrepancyRecord
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		jobs:  make(map[uuid.UUID]*model.ConsistencyJob),
		discs: make(map[uuid.UUID]*model.DiscrepancyRecord),
	}
}

func (m *mockRepo) CreateJob(ctx context.Context, job *model.ConsistencyJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) UpdateJob(ctx context.Context, job *model.ConsistencyJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*model.ConsistencyJob, error) {
	job, ok := m.jobs[id]
	if !ok {
		return nil, gormNotFound()
	}
	return job, nil
}
func (m *mockRepo) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.ConsistencyJob, error) {
	var res []model.ConsistencyJob
	for _, j := range m.jobs {
		if j.WorkspaceID == workspaceID {
			res = append(res, *j)
		}
	}
	return res, nil
}
func (m *mockRepo) CreateDiscrepancy(ctx context.Context, rec *model.DiscrepancyRecord) error {
	m.discs[rec.ID] = rec
	return nil
}
func (m *mockRepo) ListDiscrepanciesByWorkspace(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]model.DiscrepancyRecord, error) {
	var res []model.DiscrepancyRecord
	for _, d := range m.discs {
		if d.WorkspaceID == workspaceID {
			if unresolvedOnly && d.Resolved {
				continue
			}
			res = append(res, *d)
		}
	}
	return res, nil
}
func (m *mockRepo) MarkDiscrepancyResolved(ctx context.Context, id uuid.UUID) error {
	if d, ok := m.discs[id]; ok {
		d.Resolved = true
	}
	return nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string { return "record not found" }
func gormNotFound() error           { return notFoundErr{} }

func TestService_ScheduleAndRecordDiscrepancy(t *testing.T) {
	repo := newMockRepo()
	eventsBus := events.NewDispatcher(nil)
	svc := NewService(repo, nil, eventsBus, nil)

	wsID := uuid.New()
	bucketID := uuid.New()

	req := ScheduleCheckRequest{
		WorkspaceID: wsID,
		BucketID:    bucketID,
		ProviderID:  "aws-s3-east",
	}

	job, err := svc.ScheduleVerification(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected schedule error: %v", err)
	}

	err = svc.RecordDiscrepancy(context.Background(), job.ID, wsID, bucketID, uuid.New(), "archive/doc.pdf", model.DiscrepancyChecksumMismatch, "sha256 mismatch")
	if err != nil {
		t.Fatalf("failed to record discrepancy: %v", err)
	}

	err = svc.MarkJobCompleted(context.Background(), job.ID, 100, 1)
	if err != nil {
		t.Fatalf("failed to mark job completed: %v", err)
	}

	updated, _ := svc.ListJobs(context.Background(), wsID)
	if len(updated) != 1 || updated[0].Status != model.StatusFailed {
		t.Fatalf("expected job status FAILED due to discrepancy, got %+v", updated)
	}

	discs, _ := svc.ListDiscrepancies(context.Background(), wsID, true)
	if len(discs) != 1 {
		t.Fatalf("expected 1 unresolved discrepancy, got %d", len(discs))
	}
}
