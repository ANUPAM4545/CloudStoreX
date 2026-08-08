package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// EventType defines the topic name for a reliability event.
type EventType string

const (
	EventObjectUploaded         EventType = "ObjectUploaded"
	EventObjectDeleted          EventType = "ObjectDeleted"
	EventReplicationStarted     EventType = "ReplicationStarted"
	EventReplicationCompleted   EventType = "ReplicationCompleted"
	EventReplicationFailed      EventType = "ReplicationFailed"
	EventBackupCompleted        EventType = "BackupCompleted"
	EventConsistencyCheckFailed EventType = "ConsistencyCheckFailed"
	EventRepairCompleted        EventType = "RepairCompleted"
	EventFailoverStarted        EventType = "FailoverStarted"
	EventFailoverCompleted      EventType = "FailoverCompleted"
	EventRecoveryCompleted      EventType = "RecoveryCompleted"
)

// Event represents a decoupled event within the Reliability Event Bus.
type Event struct {
	ID          uuid.UUID              `json:"id"`
	Type        EventType              `json:"type"`
	WorkspaceID uuid.UUID              `json:"workspace_id,omitempty"`
	ObjectID    uuid.UUID              `json:"object_id,omitempty"`
	ProviderID  string                 `json:"provider_id,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewEvent creates a new reliability event with generated ID and current UTC timestamp.
func NewEvent(eventType EventType, workspaceID, objectID uuid.UUID, providerID string, payload map[string]interface{}) Event {
	return Event{
		ID:          uuid.New(),
		Type:        eventType,
		WorkspaceID: workspaceID,
		ObjectID:    objectID,
		ProviderID:  providerID,
		Payload:     payload,
		Timestamp:   time.Now().UTC(),
	}
}

// Handler defines a callback function executed when an event occurs.
type Handler func(ctx context.Context, event Event) error
