package events

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
)

func TestEventDispatcher_PublishSubscribe(t *testing.T) {
	dispatcher := NewDispatcher(nil)
	var handledCount int32

	dispatcher.Subscribe(EventReplicationStarted, func(ctx context.Context, event Event) error {
		atomic.AddInt32(&handledCount, 1)
		if event.Type != EventReplicationStarted {
			t.Errorf("expected %s, got %s", EventReplicationStarted, event.Type)
		}
		return nil
	})

	dispatcher.Subscribe(EventReplicationStarted, func(ctx context.Context, event Event) error {
		atomic.AddInt32(&handledCount, 1)
		return nil
	})

	evt := NewEvent(EventReplicationStarted, uuid.New(), uuid.New(), "aws-s3", nil)
	err := dispatcher.Publish(context.Background(), evt)
	if err != nil {
		t.Fatalf("unexpected error publishing event: %v", err)
	}

	if handledCount != 2 {
		t.Fatalf("expected 2 handlers to be called, got %d", handledCount)
	}
}

func TestEventDispatcher_HandlerErrorPropagation(t *testing.T) {
	dispatcher := NewDispatcher(nil)

	expectedErr := errors.New("handler failure")
	dispatcher.Subscribe(EventConsistencyCheckFailed, func(ctx context.Context, event Event) error {
		return expectedErr
	})

	evt := NewEvent(EventConsistencyCheckFailed, uuid.New(), uuid.New(), "minio", nil)
	err := dispatcher.Publish(context.Background(), evt)
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
