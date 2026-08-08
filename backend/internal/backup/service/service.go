package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/backup/model"
	"github.com/cloudstorex/backend/internal/backup/repository"
	"github.com/cloudstorex/backend/internal/jobs"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type ScheduleBackupRequest struct {
	WorkspaceID    uuid.UUID
	BucketID       uuid.UUID
	SourceProvider string
	TargetProvider string
	BackupType     model.BackupType
}

// Service coordinates backup job creation, vault record persistence, and event emissions.
type Service interface {
	ScheduleBackup(ctx context.Context, req ScheduleBackupRequest) (*model.BackupJob, error)
	MarkJobCompleted(ctx context.Context, jobID uuid.UUID, snapshotID string, objectsBackedUp, bytesBackedUp int64, vaultLoc, checksum string) error
	MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error
	ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupJob, error)
	ListRecords(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupRecord, error)
}

type defaultService struct {
	repo      repository.Repository
	jobClient jobs.Client
	eventsBus events.Publisher
	logger    *slog.Logger
}

// NewService creates a new Phase 5 Backup Engine service.
func NewService(repo repository.Repository, jobClient jobs.Client, eventsBus events.Publisher, logger *slog.Logger) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &defaultService{
		repo:      repo,
		jobClient: jobClient,
		eventsBus: eventsBus,
		logger:    logger,
	}
}

func (s *defaultService) ScheduleBackup(ctx context.Context, req ScheduleBackupRequest) (*model.BackupJob, error) {
	backupType := req.BackupType
	if backupType == "" {
		backupType = model.BackupTypeFull
	}

	job := &model.BackupJob{
		ID:             uuid.New(),
		WorkspaceID:    req.WorkspaceID,
		BucketID:       req.BucketID,
		SourceProvider: req.SourceProvider,
		TargetProvider: req.TargetProvider,
		BackupType:     backupType,
		Status:         model.BackupStatusPending,
	}

	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create backup job: %w", err)
	}

	if s.jobClient != nil {
		payload := map[string]interface{}{
			"backup_job_id":   job.ID.String(),
			"workspace_id":    req.WorkspaceID.String(),
			"bucket_id":       req.BucketID.String(),
			"source_provider": req.SourceProvider,
			"target_provider": req.TargetProvider,
			"backup_type":     string(backupType),
		}
		_, err := s.jobClient.Enqueue(ctx, "BACKUP_JOB", payload, 3)
		if err != nil {
			s.logger.Error("failed to enqueue BACKUP_JOB", slog.String("error", err.Error()))
		}
	}

	s.logger.Info("scheduled backup job",
		slog.String("job_id", job.ID.String()),
		slog.String("workspace_id", req.WorkspaceID.String()),
		slog.String("source_provider", req.SourceProvider),
		slog.String("target_provider", req.TargetProvider),
	)
	return job, nil
}

func (s *defaultService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, snapshotID string, objectsBackedUp, bytesBackedUp int64, vaultLoc, checksum string) error {
	job, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = model.BackupStatusCompleted
	job.SnapshotID = snapshotID
	job.ObjectsBackedUp = objectsBackedUp
	job.BytesBackedUp = bytesBackedUp
	job.CompletedAt = &now
	job.ErrorMessage = ""

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
	}

	record := &model.BackupRecord{
		ID:            uuid.New(),
		BackupJobID:   jobID,
		WorkspaceID:   job.WorkspaceID,
		BucketID:      job.BucketID,
		SnapshotID:    snapshotID,
		VaultLocation: vaultLoc,
		TotalBytes:    bytesBackedUp,
		Checksum:      checksum,
	}
	if err := s.repo.CreateRecord(ctx, record); err != nil {
		s.logger.Error("failed to create backup record", slog.String("error", err.Error()))
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventBackupCompleted, job.WorkspaceID, uuid.Nil, job.TargetProvider, map[string]interface{}{
			"backup_job_id":     job.ID.String(),
			"snapshot_id":       snapshotID,
			"objects_backed_up": objectsBackedUp,
			"bytes_backed_up":   bytesBackedUp,
		})
		_ = s.eventsBus.Publish(ctx, evt)
	}
	return nil
}

func (s *defaultService) MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	job, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = model.BackupStatusFailed
	job.CompletedAt = &now
	job.ErrorMessage = errMsg

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
	}
	return nil
}

func (s *defaultService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupJob, error) {
	return s.repo.ListJobsByWorkspace(ctx, workspaceID)
}

func (s *defaultService) ListRecords(ctx context.Context, workspaceID uuid.UUID) ([]model.BackupRecord, error) {
	return s.repo.ListRecordsByWorkspace(ctx, workspaceID)
}
