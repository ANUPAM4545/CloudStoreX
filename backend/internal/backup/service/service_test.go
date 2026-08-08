package service

import (
	"context"
	"testing"

	"github.com/cloudstorex/backend/internal/backup/model"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type mockRepo struct {
	jobs    map[uuid.UUID]*model.BackupJob
	records map[uuid.UUID]*model.BackupRecord
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		jobs:    make(map[uuid.UUID]*model.BackupJob),
		records: make(map[uuid.UUID]*model.BackupRecord),
	}
}

func (m *mockRepo) CreateJob(ctx context.Context, job *model.BackupJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) UpdateJob(ctx context.Context, job *model.BackupJob) error {
	m.jobs[job.ID] = job
	return nil
}
func (m *mockRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*model.BackupJob, error) {
	job, ok := m.jobs[id]
	if !ok {
		return nil, gormNotFound()
	}
	return job, nil
}
func (m *mockRepo) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupJob, error) {
	var res []model.BackupJob
	for _, j := range m.jobs {
		if j.WorkspaceID == workspaceID {
			res = append(res, *j)
		}
	}
	return res, nil
}
func (m *mockRepo) CreateRecord(ctx context.Context, rec *model.BackupRecord) error {
	m.records[rec.ID] = rec
	return nil
}
func (m *mockRepo) ListRecordsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupRecord, error) {
	var res []model.BackupRecord
	for _, r := range m.records {
		if r.WorkspaceID == workspaceID {
			res = append(res, *r)
		}
	}
	return res, nil
}

type notFoundErr struct{}

func (e notFoundErr) Error() string { return "record not found" }
func gormNotFound() error           { return notFoundErr{} }

func TestService_ScheduleAndCompleteBackup(t *testing.T) {
	repo := newMockRepo()
	eventsBus := events.NewDispatcher(nil)
	svc := NewService(repo, nil, eventsBus, nil)

	wsID := uuid.New()
	bucketID := uuid.New()

	req := ScheduleBackupRequest{
		WorkspaceID:    wsID,
		BucketID:       bucketID,
		SourceProvider: "aws-s3-east",
		TargetProvider: "aws-glacier-vault",
		BackupType:     model.BackupTypeFull,
	}

	job, err := svc.ScheduleBackup(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected schedule error: %v", err)
	}

	err = svc.MarkJobCompleted(context.Background(), job.ID, "snap-001", 100, 1048576, "vault://aws-glacier/snap-001", "sha256-mock")
	if err != nil {
		t.Fatalf("unexpected mark completed error: %v", err)
	}

	records, _ := svc.ListRecords(context.Background(), wsID)
	if len(records) != 1 || records[0].SnapshotID != "snap-001" {
		t.Fatalf("expected 1 record with snapshot snap-001, got %+v", records)
	}
}
