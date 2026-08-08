package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/chaos/model"
	"github.com/cloudstorex/backend/internal/chaos/repository"
	"github.com/cloudstorex/backend/internal/cluster"
	"github.com/google/uuid"
)

type RunExperimentRequest struct {
	WorkspaceID     uuid.UUID
	Name            string
	ExperimentType  model.ExperimentType
	TargetProvider  string
	DurationSeconds int
	TargetRPO       int
	TargetRTO       int
}

// Service manages chaos experiment lifecycles, fault injections, and RTO/RPO SLA verification.
type Service interface {
	RunExperiment(ctx context.Context, req RunExperimentRequest) (*model.ChaosExperiment, error)
	CompleteExperiment(ctx context.Context, expID uuid.UUID, measuredRPO, measuredRTO int) error
	AbortExperiment(ctx context.Context, expID uuid.UUID, reason string) error
	GetExperiment(ctx context.Context, expID uuid.UUID) (*model.ChaosExperiment, error)
	ListExperiments(ctx context.Context, workspaceID uuid.UUID) ([]model.ChaosExperiment, error)
}

type defaultService struct {
	repo     repository.Repository
	stateMgr cluster.StateManager
	logger   *slog.Logger
}

// NewService creates a new Phase 8 Chaos Engineering service.
func NewService(repo repository.Repository, stateMgr cluster.StateManager, logger *slog.Logger) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &defaultService{
		repo:     repo,
		stateMgr: stateMgr,
		logger:   logger,
	}
}

func (s *defaultService) RunExperiment(ctx context.Context, req RunExperimentRequest) (*model.ChaosExperiment, error) {
	now := time.Now().UTC()
	dur := req.DurationSeconds
	if dur <= 0 {
		dur = 60
	}

	exp := &model.ChaosExperiment{
		ID:              uuid.New(),
		WorkspaceID:     req.WorkspaceID,
		Name:            req.Name,
		ExperimentType:  req.ExperimentType,
		TargetProvider:  req.TargetProvider,
		DurationSeconds: dur,
		TargetRPO:       req.TargetRPO,
		TargetRTO:       req.TargetRTO,
		Status:          model.StatusRunning,
		StartedAt:       &now,
	}

	if err := s.repo.CreateExperiment(ctx, exp); err != nil {
		return nil, fmt.Errorf("failed to create chaos experiment: %w", err)
	}

	// Apply fault injection in Cluster State Manager for provider-level faults
	if s.stateMgr != nil && req.TargetProvider != "" {
		if req.ExperimentType == model.TypeProviderOutage {
			s.stateMgr.UpdateProviderHealth(req.TargetProvider, cluster.StateUnhealthy, 0, 10000, 1.0)
			s.logger.Warn("injected provider outage in cluster state",
				slog.String("experiment_id", exp.ID.String()),
				slog.String("provider_id", req.TargetProvider),
			)
		} else if req.ExperimentType == model.TypeLatencyInjection {
			s.stateMgr.UpdateProviderHealth(req.TargetProvider, cluster.StateDegraded, 50, 4500, 0.05)
			s.logger.Warn("injected latency degradation in cluster state",
				slog.String("experiment_id", exp.ID.String()),
				slog.String("provider_id", req.TargetProvider),
			)
		}
	}

	s.logger.Info("started chaos experiment",
		slog.String("experiment_id", exp.ID.String()),
		slog.String("name", exp.Name),
		slog.String("experiment_type", string(exp.ExperimentType)),
	)
	return exp, nil
}

func (s *defaultService) CompleteExperiment(ctx context.Context, expID uuid.UUID, measuredRPO, measuredRTO int) error {
	exp, err := s.repo.GetExperimentByID(ctx, expID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	exp.MeasuredRPO = measuredRPO
	exp.MeasuredRTO = measuredRTO
	exp.CompletedAt = &now
	exp.Status = model.StatusCompleted

	passed := true
	if exp.TargetRPO > 0 && measuredRPO > exp.TargetRPO {
		passed = false
	}
	if exp.TargetRTO > 0 && measuredRTO > exp.TargetRTO {
		passed = false
	}
	exp.Passed = passed

	// Restore normal provider health if applicable
	if s.stateMgr != nil && exp.TargetProvider != "" {
		s.stateMgr.UpdateProviderHealth(exp.TargetProvider, cluster.StateHealthy, 100, 80, 0.0)
	}

	if err := s.repo.UpdateExperiment(ctx, exp); err != nil {
		return err
	}

	s.logger.Info("completed chaos experiment",
		slog.String("experiment_id", exp.ID.String()),
		slog.Bool("passed_sla", exp.Passed),
		slog.Int("measured_rpo", measuredRPO),
		slog.Int("measured_rto", measuredRTO),
	)
	return nil
}

func (s *defaultService) AbortExperiment(ctx context.Context, expID uuid.UUID, reason string) error {
	exp, err := s.repo.GetExperimentByID(ctx, expID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	exp.Status = model.StatusAborted
	exp.CompletedAt = &now
	exp.ErrorMessage = reason
	exp.Passed = false

	if s.stateMgr != nil && exp.TargetProvider != "" {
		s.stateMgr.UpdateProviderHealth(exp.TargetProvider, cluster.StateHealthy, 100, 80, 0.0)
	}

	return s.repo.UpdateExperiment(ctx, exp)
}

func (s *defaultService) GetExperiment(ctx context.Context, expID uuid.UUID) (*model.ChaosExperiment, error) {
	return s.repo.GetExperimentByID(ctx, expID)
}

func (s *defaultService) ListExperiments(ctx context.Context, workspaceID uuid.UUID) ([]model.ChaosExperiment, error) {
	return s.repo.ListExperimentsByWorkspace(ctx, workspaceID)
}
