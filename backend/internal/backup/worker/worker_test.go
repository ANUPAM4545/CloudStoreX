package worker

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cloudstorex/backend/internal/backup/model"
	"github.com/cloudstorex/backend/internal/backup/service"
	jobModel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/google/uuid"
)

type dummyBackupService struct {
	completed bool
}

func (d *dummyBackupService) ScheduleBackup(ctx context.Context, req service.ScheduleBackupRequest) (*model.BackupJob, error) {
	return nil, nil
}
func (d *dummyBackupService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, snapshotID string, objectsBackedUp, bytesBackedUp int64, vaultLoc, checksum string) error {
	d.completed = true
	return nil
}
func (d *dummyBackupService) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return nil
}
func (d *dummyBackupService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupJob, error) {
	return nil, nil
}
func (d *dummyBackupService) ListRecords(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupRecord, error) {
	return nil, nil
}

func TestBackupWorker_Handle(t *testing.T) {
	svc := &dummyBackupService{}
	w := NewWorker(svc, nil, nil)

	payload := BackupJobPayload{
		BackupJobID:    uuid.NewString(),
		WorkspaceID:    uuid.NewString(),
		BucketID:       uuid.NewString(),
		SourceProvider: "aws-s3-east",
		TargetProvider: "aws-glacier-vault",
		BackupType:     "FULL",
	}
	data, _ := json.Marshal(payload)

	job := &jobModel.Job{
		ID:      uuid.NewString(),
		Type:    "BACKUP_JOB",
		Payload: data,
	}

	err := w.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("unexpected worker error: %v", err)
	}
	if !svc.completed {
		t.Fatalf("expected backup job to be marked completed")
	}
}
