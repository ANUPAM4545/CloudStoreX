package events_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/events"
	"github.com/stretchr/testify/assert"
)

func TestLogPublisher_Publish(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pub := events.NewLogPublisher(logger)

	event := events.Event{
		ID:        "evt-123",
		Type:      events.EventObjectCreated,
		Timestamp: time.Now(),
		ObjectID:  "obj-123",
		BucketID:  "bkt-123",
	}

	err := pub.Publish(context.Background(), event)
	assert.NoError(t, err)
}

func TestNewLogPublisher(t *testing.T) {
	pub := events.NewLogPublisher(nil)
	assert.NotNil(t, pub)
}
