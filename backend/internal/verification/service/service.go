package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/jobs"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/cloudstorex/backend/internal/verification/model"
	"github.com/cloudstorex/backend/internal/verification/repository"
	"github.com/google/uuid"
)

type ScheduleCheckRequest struct {
	WorkspaceID uuid.UUID
	BucketID    uuid.UUID
	ProviderID  string
}

// Service manages consistency verification jobs, discrepancy recording, and reliability event emissions.
type Service interface {
	ScheduleVerification(ctx context.Context, req ScheduleCheckRequest) (*model.ConsistencyJob, error)
	RecordDiscrepancy(ctx context.Context, jobID, wsID, bucketID, objID uuid.UUID, objKey string, discType model.DiscrepancyType, details string) error
	MarkJobCompleted(ctx context.Context, jobID uuid.UUID, totalChecked, discrepanciesFound int64) error
	MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error
	ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.ConsistencyJob, error)
	ListDiscrepancies(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]model.DiscrepancyRecord, error)
}

type defaultService struct {
	repo      repository.Repository
	jobClient jobs.Client
	eventsBus events.Publisher
	logger    *slog.Logger
}

// NewService creates a new Phase 6 Consistency Verification service.
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

func (s *defaultService) ScheduleVerification(ctx context.Context, req ScheduleCheckRequest) (*model.ConsistencyJob, error) {
	job := &model.ConsistencyJob{
		ID:          uuid.New(),
		WorkspaceID: req.WorkspaceID,
		BucketID:    req.BucketID,
		ProviderID:  req.ProviderID,
		Status:      model.StatusPending,
	}

	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create consistency check job: %w", err)
	}

	if s.jobClient != nil {
		payload := map[string]interface{}{
			"job_id":       job.ID.String(),
			"workspace_id": req.WorkspaceID.String(),
			"bucket_id":    req.BucketID.String(),
			"provider_id":  req.ProviderID,
		}
		_, err := s.jobClient.Enqueue(ctx, "CONSISTENCY_CHECK", payload, 3)
		if err != nil {
			s.logger.Error("failed to enqueue CONSISTENCY_CHECK", slog.String("error", err.Error()))
		}
	}

	s.logger.Info("scheduled consistency verification job",
		slog.String("job_id", job.ID.String()),
		slog.String("workspace_id", req.WorkspaceID.String()),
		slog.String("provider_id", req.ProviderID),
	)
	return job, nil
}

func (s *defaultService) RecordDiscrepancy(ctx context.Context, jobID, wsID, bucketID, objID uuid.UUID, objKey string, discType model.DiscrepancyType, details string) error {
	rec := &model.DiscrepancyRecord{
		ID:          uuid.New(),
		JobID:       jobID,
		WorkspaceID: wsID,
		BucketID:    bucketID,
		ObjectID:    objID,
		ObjectKey:   objKey,
		Type:        discType,
		Details:     details,
		Resolved:    false,
	}
	if err := s.repo.CreateDiscrepancy(ctx, rec); err != nil {
		return err
	}

	s.logger.Warn("recorded object consistency discrepancy",
		slog.String("job_id", jobID.String()),
		slog.String("object_key", objKey),
		slog.String("discrepancy_type", string(discType)),
		slog.String("details", details),
	)

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventConsistencyCheckFailed, wsID, objID, "", map[string]interface{}{
			"job_id":           jobID.String(),
			"bucket_id":        bucketID.String(),
			"object_key":       objKey,
			"discrepancy_type": string(discType),
			"details":          details,
		})
		_ = s.eventsBus.Publish(ctx, evt)
	}
	return nil
}

func (s *defaultService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, totalChecked, discrepanciesFound int64) error {
	job, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = model.StatusPassed
	if discrepanciesFound > 0 {
		job.Status = model.StatusFailed
	}
	job.TotalObjectsChecked = totalChecked
	job.DiscrepanciesFound = discrepanciesFound
	job.CompletedAt = &now
	job.ErrorMessage = ""

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
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

func (s *defaultService) ListJobs(ctx context.Context, workspaceID uuid.UUID) ([]model.ConsistencyJob, error) {
	return s.repo.ListJobsByWorkspace(ctx, workspaceID)
}

func (s *defaultService) ListDiscrepancies(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]model.DiscrepancyRecord, error) {
	return s.repo.ListDiscrepanciesByWorkspace(ctx, workspaceID, unresolvedOnly)
}
