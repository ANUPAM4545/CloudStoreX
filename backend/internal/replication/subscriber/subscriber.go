package subscriber

import (
	"context"
	"log/slog"

	"github.com/cloudstorex/backend/internal/replication/model"
	"github.com/cloudstorex/backend/internal/replication/service"
	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/google/uuid"
)

// RegisterReplicationSubscriber subscribes the replication service to ObjectUploaded events (Epic 13 Phase 2).
func RegisterReplicationSubscriber(dispatcher events.Subscriber, replSvc service.Service, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	handler := func(ctx context.Context, evt events.Event) error {
		logger.Info("received ObjectUploaded reliability event for replication",
			slog.String("workspace_id", evt.WorkspaceID.String()),
			slog.String("object_id", evt.ObjectID.String()),
		)

		payload := evt.Payload
		if payload == nil {
			return nil
		}

		bucketIDStr, _ := payload["bucket_id"].(string)
		bucketID, _ := uuid.Parse(bucketIDStr)
		objectKey, _ := payload["object_key"].(string)
		primaryProvider, _ := payload["primary_provider"].(string)
		modeStr, _ := payload["mode"].(string)
		mode := model.ReplicationMode(modeStr)
		if mode == "" {
			mode = model.ModeAsync
		}

		var replicaProviders []string
		if reps, ok := payload["replica_providers"].([]string); ok {
			replicaProviders = reps
		} else if repsInterface, ok := payload["replica_providers"].([]interface{}); ok {
			for _, r := range repsInterface {
				if s, ok := r.(string); ok {
					replicaProviders = append(replicaProviders, s)
				}
			}
		}

		for _, repProvider := range replicaProviders {
			if repProvider == "" || repProvider == primaryProvider {
				continue
			}
			req := service.ScheduleRequest{
				WorkspaceID:     evt.WorkspaceID,
				BucketID:        bucketID,
				ObjectID:        evt.ObjectID,
				ObjectKey:       objectKey,
				PrimaryProvider: primaryProvider,
				ReplicaProvider: repProvider,
				Mode:            mode,
			}
			_, err := replSvc.ScheduleReplication(ctx, req)
			if err != nil {
				logger.Error("failed to schedule replication from event",
					slog.String("replica_provider", repProvider),
					slog.String("error", err.Error()),
				)
			}
		}
		return nil
	}

	dispatcher.Subscribe(events.EventObjectUploaded, handler)
	return nil
}
