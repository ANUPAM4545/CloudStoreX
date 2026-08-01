package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/jobs"
)

// JobEventPublisher implements Publisher by enqueuing jobs into the background job system.
type JobEventPublisher struct {
	jobClient jobs.Client
	log       *slog.Logger
}

func NewJobEventPublisher(jobClient jobs.Client, log *slog.Logger) *JobEventPublisher {
	if log == nil {
		log = slog.Default()
	}
	return &JobEventPublisher{
		jobClient: jobClient,
		log:       log,
	}
}

func (p *JobEventPublisher) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Determine job type and delay based on event type
	jobType := "ProcessEvent"
	var executeAt *time.Time

	switch event.Type {
	case EventObjectDeleted:
		// For soft delete, we might schedule a hard cleanup job in 30 days.
		// For Epic 10 Phase 1, we can just log or enqueue it. Let's schedule it 30 days from now (or a shorter time for testing).
		jobType = "SoftDeleteCleanup"
		future := time.Now().Add(30 * 24 * time.Hour)
		executeAt = &future
	}

	var jobID string
	if executeAt != nil {
		job, err := p.jobClient.EnqueueScheduled(ctx, jobType, payload, 3, *executeAt)
		if err != nil {
			p.log.ErrorContext(ctx, "failed to enqueue job for event",
				slog.String("event_type", string(event.Type)),
				slog.String("error", err.Error()),
			)
			return err
		}
		jobID = job.ID
	} else {
		job, err := p.jobClient.Enqueue(ctx, jobType, payload, 3)
		if err != nil {
			p.log.ErrorContext(ctx, "failed to enqueue job for event",
				slog.String("event_type", string(event.Type)),
				slog.String("error", err.Error()),
			)
			return err
		}
		jobID = job.ID
	}

	p.log.InfoContext(ctx, "enqueued job for event",
		slog.String("event_type", string(event.Type)),
		slog.String("job_id", jobID),
	)

	return nil
}
