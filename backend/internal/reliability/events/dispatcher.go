package events

import (
	"context"
	"log/slog"
	"sync"
)

// Publisher defines the interface for emitting reliability events.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

// Subscriber defines the interface for subscribing to reliability events.
type Subscriber interface {
	Subscribe(eventType EventType, handler Handler)
}

// Dispatcher coordinates event publishing and subscriber delivery across reliability subsystems.
type Dispatcher interface {
	Publisher
	Subscriber
}

type defaultDispatcher struct {
	mu        sync.RWMutex
	handlers  map[EventType][]Handler
	logger    *slog.Logger
	asyncPool chan Event
}

// NewDispatcher creates a new thread-safe Reliability Event Bus dispatcher.
func NewDispatcher(logger *slog.Logger) Dispatcher {
	if logger == nil {
		logger = slog.Default()
	}
	return &defaultDispatcher{
		handlers:  make(map[EventType][]Handler),
		logger:    logger,
		asyncPool: make(chan Event, 1000),
	}
}

// Subscribe registers an event handler for a specific topic.
func (d *defaultDispatcher) Subscribe(eventType EventType, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Publish delivers an event to all subscribed handlers for its EventType.
func (d *defaultDispatcher) Publish(ctx context.Context, event Event) error {
	d.mu.RLock()
	handlers := d.handlers[event.Type]
	d.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			d.logger.Error("failed to handle reliability event",
				slog.String("event_type", string(event.Type)),
				slog.String("event_id", event.ID.String()),
				slog.String("error", err.Error()),
			)
			return err
		}
	}
	return nil
}
