package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/jobs"
	"github.com/cloudstorex/backend/internal/disasterrecovery/model"
	"github.com/cloudstorex/backend/internal/disasterrecovery/repository"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type CreatePlanRequest struct {
	WorkspaceID    uuid.UUID
	Name           string
	Description    string
	SourceProvider string
	TargetProvider string
	RPOSeconds     int
	RTOSeconds     int
}

type TriggerRecoveryRequest struct {
	PlanID      uuid.UUID
	PointInTime *time.Time
}

// Service manages disaster recovery plans, point-in-time restore jobs, and event emissions.
type Service interface {
	CreatePlan(ctx context.Context, req CreatePlanRequest) (*model.RecoveryPlan, error)
	GetPlan(ctx context.Context, id uuid.UUID) (*model.RecoveryPlan, error)
	ListPlans(ctx context.Context, workspaceID uuid.UUID) ([]model.RecoveryPlan, error)
	TriggerRecovery(ctx context.Context, req TriggerRecoveryRequest) (*model.RecoveryJob, error)
	MarkJobCompleted(ctx context.Context, jobID uuid.UUID, objectsRestored, bytesRestored int64) error
	MarkJobFailed(ctx context.Context, jobID uuid.UUID, errMsg string) error
	GetJob(ctx context.Context, id uuid.UUID) (*model.RecoveryJob, error)
}

type defaultService struct {
	repo      repository.Repository
	jobClient jobs.Client
	eventsBus events.Publisher
	logger    *slog.Logger
}

// NewService creates a new Phase 3 Disaster Recovery service.
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

func (s *defaultService) CreatePlan(ctx context.Context, req CreatePlanRequest) (*model.RecoveryPlan, error) {
	plan := &model.RecoveryPlan{
		ID:             uuid.New(),
		WorkspaceID:    req.WorkspaceID,
		Name:           req.Name,
		Description:    req.Description,
		SourceProvider: req.SourceProvider,
		TargetProvider: req.TargetProvider,
		RPOSeconds:     req.RPOSeconds,
		RTOSeconds:     req.RTOSeconds,
		Status:         model.PlanStatusActive,
	}

	if plan.RPOSeconds <= 0 {
		plan.RPOSeconds = 3600
	}
	if plan.RTOSeconds <= 0 {
		plan.RTOSeconds = 1800
	}

	if err := s.repo.CreatePlan(ctx, plan); err != nil {
		return nil, fmt.Errorf("failed to create recovery plan: %w", err)
	}

	s.logger.Info("created disaster recovery plan",
		slog.String("plan_id", plan.ID.String()),
		slog.String("workspace_id", plan.WorkspaceID.String()),
		slog.String("source_provider", plan.SourceProvider),
		slog.String("target_provider", plan.TargetProvider),
	)
	return plan, nil
}

func (s *defaultService) GetPlan(ctx context.Context, id uuid.UUID) (*model.RecoveryPlan, error) {
	return s.repo.GetPlanByID(ctx, id)
}

func (s *defaultService) ListPlans(ctx context.Context, workspaceID uuid.UUID) ([]model.RecoveryPlan, error) {
	return s.repo.ListPlansByWorkspace(ctx, workspaceID)
}

func (s *defaultService) TriggerRecovery(ctx context.Context, req TriggerRecoveryRequest) (*model.RecoveryJob, error) {
	plan, err := s.repo.GetPlanByID(ctx, req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("failed to load recovery plan: %w", err)
	}

	job := &model.RecoveryJob{
		ID:             uuid.New(),
		PlanID:         plan.ID,
		WorkspaceID:    plan.WorkspaceID,
		SourceProvider: plan.SourceProvider,
		TargetProvider: plan.TargetProvider,
		PointInTime:    req.PointInTime,
		Status:         model.JobStatusPending,
	}

	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create recovery job: %w", err)
	}

	if s.jobClient != nil {
		payload := map[string]interface{}{
			"job_id":          job.ID.String(),
			"plan_id":         plan.ID.String(),
			"workspace_id":    plan.WorkspaceID.String(),
			"source_provider": plan.SourceProvider,
			"target_provider": plan.TargetProvider,
			"point_in_time":   req.PointInTime,
		}
		_, err := s.jobClient.Enqueue(ctx, "RESTORE_JOB", payload, 3)
		if err != nil {
			s.logger.Error("failed to enqueue RESTORE_JOB", slog.String("error", err.Error()))
		}
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventRecoveryCompleted, plan.WorkspaceID, uuid.Nil, plan.TargetProvider, map[string]interface{}{
			"job_id": job.ID.String(),
			"status": "TRIGGERED",
		})
		_ = s.eventsBus.Publish(ctx, evt)
	}

	s.logger.Info("triggered disaster recovery restore job",
		slog.String("job_id", job.ID.String()),
		slog.String("plan_id", plan.ID.String()),
		slog.String("source_provider", plan.SourceProvider),
		slog.String("target_provider", plan.TargetProvider),
	)
	return job, nil
}

func (s *defaultService) MarkJobCompleted(ctx context.Context, jobID uuid.UUID, objectsRestored, bytesRestored int64) error {
	job, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	job.Status = model.JobStatusCompleted
	job.ObjectsRestored = objectsRestored
	job.BytesRestored = bytesRestored
	job.CompletedAt = &now
	job.ErrorMessage = ""

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventRecoveryCompleted, job.WorkspaceID, uuid.Nil, job.TargetProvider, map[string]interface{}{
			"job_id":           job.ID.String(),
			"objects_restored": objectsRestored,
			"bytes_restored":   bytesRestored,
			"status":           "COMPLETED",
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
	job.Status = model.JobStatusFailed
	job.CompletedAt = &now
	job.ErrorMessage = errMsg

	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return err
	}
	return nil
}

func (s *defaultService) GetJob(ctx context.Context, id uuid.UUID) (*model.RecoveryJob, error) {
	return s.repo.GetJobByID(ctx, id)
}
