package events

import (
	"context"
	"log/slog"

	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/cloudstorex/backend/internal/policy/repository"
)

// Publisher defines an asynchronous event publisher for routing decisions.
type Publisher interface {
	PublishRoutingDecision(ctx context.Context, decision *model.RoutingDecision) error
}

type logPublisher struct {
	logger *slog.Logger
	repo   repository.Repository
}

// NewLogPublisher creates a simple log-based event publisher.
func NewLogPublisher(logger *slog.Logger, repo repository.Repository) Publisher {
	if logger == nil {
		logger = slog.Default()
	}
	return &logPublisher{logger: logger, repo: repo}
}

func (p *logPublisher) PublishRoutingDecision(ctx context.Context, decision *model.RoutingDecision) error {
	// Async event emission
	go func() {
		p.logger.Info("RoutingDecisionCreated",
			slog.String("workspace_id", decision.WorkspaceID.String()),
			slog.String("object_key", decision.ObjectKey),
			slog.String("provider_id", decision.ProviderID),
			slog.Bool("is_fallback", decision.IsFallback),
			slog.String("matched_rule", decision.MatchedRule),
			slog.String("explanation", decision.DecisionReason),
			slog.Int64("latency_ms", decision.LatencyMs),
		)
		
		if p.repo != nil {
			if err := p.repo.CreateRoutingDecision(context.Background(), decision); err != nil {
				p.logger.Error("Failed to save RoutingDecision to database", slog.String("error", err.Error()))
			}
		}
	}()
	return nil
}
