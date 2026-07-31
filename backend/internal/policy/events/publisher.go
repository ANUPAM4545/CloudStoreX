package events

import (
	"context"
	"log/slog"

	"github.com/cloudstorex/backend/internal/policy/model"
)

// Publisher defines an asynchronous event publisher for routing decisions.
type Publisher interface {
	PublishRoutingDecision(ctx context.Context, decision *model.RoutingDecision) error
}

type logPublisher struct {
	logger *slog.Logger
	// In the future, this would hold a reference to Kafka/RabbitMQ/PubSub client
}

// NewLogPublisher creates a simple log-based event publisher.
func NewLogPublisher(logger *slog.Logger) Publisher {
	if logger == nil {
		logger = slog.Default()
	}
	return &logPublisher{logger: logger}
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
	}()
	return nil
}
