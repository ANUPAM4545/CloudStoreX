package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cloudstorex/backend/internal/cluster"
	"github.com/cloudstorex/backend/internal/failover/model"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

type FailoverConfig struct {
	CooldownPeriod time.Duration
}

// Service manages automatic failover and failback using cluster state and provider weighting.
type Service interface {
	CheckAndFailover(ctx context.Context, workspaceID uuid.UUID, primaryProvider string, candidates []string) (*model.FailoverEvent, error)
	CheckAndFailback(ctx context.Context, failoverEventID uuid.UUID) (*model.FailoverEvent, error)
	GetActiveFailovers() []*model.FailoverEvent
}

type defaultService struct {
	mu             sync.RWMutex
	stateMgr       cluster.StateManager
	eventsBus      events.Publisher
	logger         *slog.Logger
	config         FailoverConfig
	activeEvents   map[uuid.UUID]*model.FailoverEvent
	lastFailoverAt map[string]time.Time // key: workspaceID+provider
}

// NewService creates a new Phase 4 Automatic Failover service.
func NewService(stateMgr cluster.StateManager, eventsBus events.Publisher, config FailoverConfig, logger *slog.Logger) Service {
	if logger == nil {
		logger = slog.Default()
	}
	if config.CooldownPeriod == 0 {
		config.CooldownPeriod = 5 * time.Minute
	}
	return &defaultService{
		stateMgr:       stateMgr,
		eventsBus:      eventsBus,
		logger:         logger,
		config:         config,
		activeEvents:   make(map[uuid.UUID]*model.FailoverEvent),
		lastFailoverAt: make(map[string]time.Time),
	}
}

func (s *defaultService) CheckAndFailover(ctx context.Context, workspaceID uuid.UUID, primaryProvider string, candidates []string) (*model.FailoverEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("%s:%s", workspaceID.String(), primaryProvider)
	if last, exists := s.lastFailoverAt[key]; exists {
		if time.Since(last) < s.config.CooldownPeriod {
			return nil, fmt.Errorf("failover cooldown in effect for %s (remaining: %v)", primaryProvider, s.config.CooldownPeriod-time.Since(last))
		}
	}

	// 1. Verify primary provider health in Cluster State Manager
	health, ok := s.stateMgr.GetProviderHealth(primaryProvider)
	if ok && health.State == cluster.StateHealthy {
		return nil, fmt.Errorf("primary provider %s is still healthy; failover not required", primaryProvider)
	}

	// 2. Select highest priority healthy candidate
	var selectedCandidate string
	for _, cand := range candidates {
		candHealth, ok := s.stateMgr.GetProviderHealth(cand)
		if ok && candHealth.State == cluster.StateHealthy {
			selectedCandidate = cand
			break
		}
	}

	if selectedCandidate == "" {
		return nil, fmt.Errorf("no healthy secondary candidate providers available for failover from %s", primaryProvider)
	}

	evt := &model.FailoverEvent{
		ID:                 uuid.New(),
		WorkspaceID:        workspaceID,
		OriginalProvider:   primaryProvider,
		NewPrimaryProvider: selectedCandidate,
		Reason:             fmt.Sprintf("primary provider %s unhealthy; failed over to candidate %s", primaryProvider, selectedCandidate),
		Status:             model.StatusActive,
		TriggeredAt:        time.Now().UTC(),
	}

	s.activeEvents[evt.ID] = evt
	s.lastFailoverAt[key] = time.Now().UTC()

	if s.eventsBus != nil {
		event := events.NewEvent(events.EventFailoverCompleted, workspaceID, uuid.Nil, selectedCandidate, map[string]interface{}{
			"failover_id":       evt.ID.String(),
			"original_provider": primaryProvider,
			"new_provider":      selectedCandidate,
		})
		_ = s.eventsBus.Publish(ctx, event)
	}

	s.logger.Warn("executed automatic provider failover",
		slog.String("failover_id", evt.ID.String()),
		slog.String("workspace_id", workspaceID.String()),
		slog.String("original_provider", primaryProvider),
		slog.String("new_primary_provider", selectedCandidate),
	)
	return evt, nil
}

func (s *defaultService) CheckAndFailback(ctx context.Context, failoverEventID uuid.UUID) (*model.FailoverEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	evt, exists := s.activeEvents[failoverEventID]
	if !exists {
		return nil, fmt.Errorf("active failover event %s not found", failoverEventID)
	}

	// Check if original provider has recovered and is healthy
	health, ok := s.stateMgr.GetProviderHealth(evt.OriginalProvider)
	if !ok || health.State != cluster.StateHealthy {
		return nil, fmt.Errorf("original provider %s is not healthy yet; cannot failback", evt.OriginalProvider)
	}

	now := time.Now().UTC()
	evt.Status = model.StatusFailback
	evt.ResolvedAt = &now
	delete(s.activeEvents, failoverEventID)

	if s.eventsBus != nil {
		event := events.NewEvent(events.EventFailoverCompleted, evt.WorkspaceID, uuid.Nil, evt.OriginalProvider, map[string]interface{}{
			"failover_id":      evt.ID.String(),
			"restored_primary": evt.OriginalProvider,
			"status":           "FAILBACK_COMPLETED",
		})
		_ = s.eventsBus.Publish(ctx, event)
	}

	s.logger.Info("completed automatic failback to original primary provider",
		slog.String("failover_id", evt.ID.String()),
		slog.String("restored_primary_provider", evt.OriginalProvider),
	)
	return evt, nil
}

func (s *defaultService) GetActiveFailovers() []*model.FailoverEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*model.FailoverEvent, 0, len(s.activeEvents))
	for _, evt := range s.activeEvents {
		res = append(res, evt)
	}
	return res
}
