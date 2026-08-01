package events

import (
	"context"
	"log/slog"
	"time"
)

// EventType defines the type of a metadata event.
type EventType string

const (
	EventObjectCreated    EventType = "ObjectCreated"
	EventObjectUpdated    EventType = "ObjectUpdated"
	EventObjectDeleted    EventType = "ObjectDeleted"
	EventObjectRestored   EventType = "ObjectRestored"
	EventObjectTagged     EventType = "ObjectTagged"
	EventObjectUntagged   EventType = "ObjectUntagged"
	EventBucketCreated    EventType = "BucketCreated"
	EventBucketDeleted    EventType = "BucketDeleted"
	EventMetadataSync     EventType = "MetadataSync"

	EventRetentionApplied  EventType = "RetentionApplied"
	EventLegalHoldEnabled  EventType = "LegalHoldEnabled"
	EventLegalHoldReleased EventType = "LegalHoldReleased"
	EventQuotaExceeded     EventType = "QuotaExceeded"
	EventLifecycleExpired  EventType = "LifecycleExpired"
)

// Event payload encapsulates the lifecycle event.
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	BucketID  string                 `json:"bucket_id,omitempty"`
	ObjectID  string                 `json:"object_id,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// Publisher defines an interface for emitting metadata events.
// For initial implementation, this just logs the event.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type LogPublisher struct {
	log *slog.Logger
}

func NewLogPublisher(log *slog.Logger) *LogPublisher {
	if log == nil {
		log = slog.Default()
	}
	return &LogPublisher{log: log}
}

func (p *LogPublisher) Publish(ctx context.Context, event Event) error {
	p.log.InfoContext(ctx, "metadata_event_published",
		slog.String("event_type", string(event.Type)),
		slog.String("bucket_id", event.BucketID),
		slog.String("object_id", event.ObjectID),
		slog.Any("payload", event.Payload),
	)
	return nil
}
