package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/jobs"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/cloudstorex/backend/internal/selfhealing/model"
	"github.com/cloudstorex/backend/internal/selfhealing/repository"
	"github.com/google/uuid"
)

type ScheduleRepairRequest struct {
	WorkspaceID    uuid.UUID
	BucketID       uuid.UUID
	ObjectID       uuid.UUID
	ObjectKey      string
	TargetProvider string
	SourceProvider string
	SourceType     model.RepairSource
}

// Service coordinates automated self-healing repair jobs and reliability events.
type Service interface {
	ScheduleRepair(ctx context.Context, req ScheduleRepairRequest) (*model.RepairJob, error)
	MarkJobCompleted(ctx context.Context, jobID uuid.UUID, bytesRepaired int64) error
	MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error
	ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error)
}

type defaultService struct {
	repo      repository.Repository
	jobClient jobs.Client
	eventsBus events.Publisher
	logger    *slog.Logger
}

// NewService creates a new Phase 7 Self-Healing service.
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

func (s *defaultService) ScheduleRepair(ctx context.Context, req ScheduleRepairRequest) (*model.RepairJob, error) {
	sourceType := req.SourceType
	if sourceType == "" {
		sourceType = model.SourceReplica
	}

	job := &model.RepairJob{
		ID:             uuid.New(),
		WorkspaceID:    req.WorkspaceID,
		BucketID:       req.BucketID,
		ObjectID:       req.ObjectID,
		ObjectKey:      req.ObjectKey,
		TargetProvider: req.TargetProvider,
		SourceProvider: req.SourceProvider,
		SourceType:     sourceType,
		Status:         model.StatusPending,
	}

	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create repair job: %w", err)
	}

	if s.jobClient != nil {
		payload := map[string]interface{}{
			"repair_job_id":   job.ID.String(),
			"workspace_id":    req.WorkspaceID.String(),
			"bucket_id":       req.BucketID.String(),
			"object_id":       req.ObjectID.String(),
			"object_key":      req.ObjectKey,
			"target_provider": req.TargetProvider,
			"source_provider": req.SourceProvider,
			"source_type":     string(sourceType),
		}
		_, err := s.jobClient.Enqueue(ctx, "SELF_HEAL_JOB", payload, 3)
		if err != nil {
			s.logger.Error("failed to enqueue SELF_HEAL_JOB", slog.String("error", err.Error()))
		}
	}

	s.logger.Info("scheduled self-healing repair job",
		slog.String("job_id", job.ID.String()),
		slog.String("object_key", req.ObjectKey),
		slog.String("source_provider", req.SourceProvider),
		slog.String("target_provider", req.TargetProvider),
	)
	return job, nil
}

func (s *defaultService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, bytesRepaired int64) error {
	job, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = model.StatusCompleted
	job.BytesRepaired = bytesRepaired
	job.CompletedAt = &now
	job.ErrorMessage = ""

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventRepairCompleted, job.WorkspaceID, job.ObjectID, job.TargetProvider, map[string]interface{}{
			"repair_job_id":  job.ID.String(),
			"object_key":     job.ObjectKey,
			"bytes_repaired": bytesRepaired,
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
	job.Status = model.StatusFailed
	job.CompletedAt = &now
	job.ErrorMessage = errMsg

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
	}
	return nil
}

func (s *defaultService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.RepairJob, error) {
	return s.repo.ListJobsByWorkspace(ctx, workspaceID)
}
