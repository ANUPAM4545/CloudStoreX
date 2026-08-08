package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/jobs"
	"github.com/cloudstorex/backend/internal/replication/model"
	"github.com/cloudstorex/backend/internal/replication/repository"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type ScheduleRequest struct {
	WorkspaceID     uuid.UUID
	BucketID        uuid.UUID
	ObjectID        uuid.UUID
	ObjectKey       string
	PrimaryProvider string
	ReplicaProvider string
	Mode            model.ReplicationMode
}

// Service manages replication scheduling, state updates, and job queue integration.
type Service interface {
	ScheduleReplication(ctx context.Context, req ScheduleRequest) (*model.ObjectReplication, error)
	MarkCompleted(ctx context.Context, id uuid.UUID, checksum string, bytesReplicated int64) error
	MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error
	GetReplicationsByObject(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error)
}

type defaultService struct {
	repo      repository.Repository
	jobClient jobs.Client
	eventsBus events.Publisher
	logger    *slog.Logger
}

// NewService creates a new Phase 1 Replication Engine service.
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

func (s *defaultService) ScheduleReplication(ctx context.Context, req ScheduleRequest) (*model.ObjectReplication, error) {
	repl := &model.ObjectReplication{
		ID:              uuid.New(),
		WorkspaceID:     req.WorkspaceID,
		BucketID:        req.BucketID,
		ObjectID:        req.ObjectID,
		ObjectKey:       req.ObjectKey,
		PrimaryProvider: req.PrimaryProvider,
		ReplicaProvider: req.ReplicaProvider,
		Mode:            req.Mode,
		Status:          model.StatusPending,
	}

	if err := s.repo.Create(ctx, repl); err != nil {
		return nil, fmt.Errorf("failed to create replication record: %w", err)
	}

	// Enqueue asynchronous job for REPLICATE_OBJECT
	if s.jobClient != nil {
		payload := map[string]interface{}{
			"replication_id":   repl.ID.String(),
			"workspace_id":     req.WorkspaceID.String(),
			"bucket_id":        req.BucketID.String(),
			"object_id":        req.ObjectID.String(),
			"object_key":       req.ObjectKey,
			"primary_provider": req.PrimaryProvider,
			"replica_provider": req.ReplicaProvider,
		}
		_, err := s.jobClient.Enqueue(ctx, "REPLICATE_OBJECT", payload, 3)
		if err != nil {
			s.logger.Error("failed to enqueue REPLICATE_OBJECT job", slog.String("error", err.Error()))
		}
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventReplicationStarted, req.WorkspaceID, req.ObjectID, req.ReplicaProvider, map[string]interface{}{
			"replication_id":   repl.ID.String(),
			"primary_provider": req.PrimaryProvider,
		})
		_ = s.eventsBus.Publish(ctx, evt)
	}

	s.logger.Info("scheduled object replication",
		slog.String("replication_id", repl.ID.String()),
		slog.String("object_id", req.ObjectID.String()),
		slog.String("primary_provider", req.PrimaryProvider),
		slog.String("replica_provider", req.ReplicaProvider),
	)
	return repl, nil
}

func (s *defaultService) MarkCompleted(ctx context.Context, id uuid.UUID, checksum string, bytesReplicated int64) error {
	repl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	repl.Status = model.StatusCompleted
	repl.Checksum = checksum
	repl.BytesReplicated = bytesReplicated
	repl.ReplicatedAt = &now
	repl.LastVerifiedAt = &now
	repl.ErrorMessage = ""

	if err := s.repo.Update(ctx, repl); err != nil {
		return err
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventReplicationCompleted, repl.WorkspaceID, repl.ObjectID, repl.ReplicaProvider, map[string]interface{}{
			"replication_id":   repl.ID.String(),
			"checksum":         checksum,
			"bytes_replicated": bytesReplicated,
		})
		_ = s.eventsBus.Publish(ctx, evt)
	}
	return nil
}

func (s *defaultService) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	repl, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	repl.Status = model.StatusFailed
	repl.ErrorMessage = errMsg

	if err := s.repo.Update(ctx, repl); err != nil {
		return err
	}

	if s.eventsBus != nil {
		evt := events.NewEvent(events.EventReplicationFailed, repl.WorkspaceID, repl.ObjectID, repl.ReplicaProvider, map[string]interface{}{
			"replication_id": repl.ID.String(),
			"error":          errMsg,
		})
		_ = s.eventsBus.Publish(ctx, evt)
	}
	return nil
}

func (s *defaultService) GetReplicationsByObject(ctx context.Context, objectID uuid.UUID) ([]model.ObjectReplication, error) {
	return s.repo.ListByObjectID(ctx, objectID)
}
