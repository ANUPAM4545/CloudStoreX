package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	defaultQueue = "cloudstorex:jobs:queue:default"
)

// Client is responsible for enqueuing jobs into the system.
type Client interface {
	Enqueue(ctx context.Context, jobType string, payload interface{}, maxRetries int) (*model.Job, error)
	EnqueueScheduled(ctx context.Context, jobType string, payload interface{}, maxRetries int, executeAt time.Time) (*model.Job, error)
}

type defaultClient struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewClient creates a new job client that persists to Postgres and enqueues via Redis.
func NewClient(db *gorm.DB, redis *redis.Client) Client {
	return &defaultClient{
		db:    db,
		redis: redis,
	}
}

func (c *defaultClient) Enqueue(ctx context.Context, jobType string, payload interface{}, maxRetries int) (*model.Job, error) {
	return c.EnqueueScheduled(ctx, jobType, payload, maxRetries, time.Now())
}

func (c *defaultClient) EnqueueScheduled(ctx context.Context, jobType string, payload interface{}, maxRetries int, executeAt time.Time) (*model.Job, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job payload: %w", err)
	}

	job := &model.Job{
		ID:          uuid.New().String(),
		Type:        jobType,
		Payload:     datatypes.JSON(payloadBytes),
		State:       model.JobStatePending,
		MaxRetries:  maxRetries,
		ScheduledAt: &executeAt,
	}

	// Persist the job to PostgreSQL for durability
	if err := c.db.WithContext(ctx).Create(job).Error; err != nil {
		return nil, fmt.Errorf("failed to persist job to db: %w", err)
	}

	jobLog := &model.JobLog{
		JobID:   uuid.MustParse(job.ID),
		State:   model.JobStatePending,
		Message: "Job enqueued",
	}
	_ = c.db.WithContext(ctx).Create(jobLog) // fire and forget for log

	// In a full production system, a scheduler would periodically scan DB for scheduled jobs
	// and push them to Redis. For now, we push immediately (if not scheduled far in future)
	// or rely on a simple ZSET in redis. For simplicity, we just use LPUSH for all.
	jobJSON, _ := json.Marshal(job)
	
	if err := c.redis.LPush(ctx, defaultQueue, jobJSON).Err(); err != nil {
		// Even if Redis fails, the job is in DB and a sweeper could pick it up.
		return job, fmt.Errorf("failed to enqueue to redis (but saved to db): %w", err)
	}

	return job, nil
}
