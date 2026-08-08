package subscriber

import (
	"context"
	"log/slog"

	"github.com/cloudstorex/backend/internal/reliability/events"
	"github.com/cloudstorex/backend/internal/selfhealing/model"
	"github.com/cloudstorex/backend/internal/selfhealing/service"
	"github.com/google/uuid"
)

// RegisterSelfHealingSubscriber subscribes the self-healing service to ConsistencyCheckFailed reliability events.
func RegisterSelfHealingSubscriber(dispatcher events.Subscriber, repairSvc service.Service, logger *slog.Logger) error {
	if logger == nil {
		logger = slog.Default()
	}

	handler := func(ctx context.Context, evt events.Event) error {
		logger.Warn("received ConsistencyCheckFailed event; triggering self-healing repair",
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

		req := service.ScheduleRepairRequest{
			WorkspaceID:    evt.WorkspaceID,
			BucketID:       bucketID,
			ObjectID:       evt.ObjectID,
			ObjectKey:      objectKey,
			TargetProvider: evt.ProviderID,
			SourceProvider: "aws-s3-west",
			SourceType:     model.SourceReplica,
		}
		if req.TargetProvider == "" {
			req.TargetProvider = "aws-s3-east"
		}

		_, err := repairSvc.ScheduleRepair(ctx, req)
		if err != nil {
			logger.Error("failed to schedule self-healing repair from consistency event",
				slog.String("object_key", objectKey),
				slog.String("error", err.Error()),
			)
		}
		return nil
	}

	dispatcher.Subscribe(events.EventConsistencyCheckFailed, handler)
	return nil
}
