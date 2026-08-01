package jobs

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Dispatcher manages the worker pool and job fetching loop.
type Dispatcher interface {
	Start(ctx context.Context)
	Stop()
}

type defaultDispatcher struct {
	db       *gorm.DB
	redis    *redis.Client
	registry Registry
	workers  int
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewDispatcher creates a job dispatcher that manages workers executing jobs.
func NewDispatcher(db *gorm.DB, redis *redis.Client, registry Registry, numWorkers int) Dispatcher {
	return &defaultDispatcher{
		db:       db,
		redis:    redis,
		registry: registry,
		workers:  numWorkers,
		stopChan: make(chan struct{}),
	}
}

func (d *defaultDispatcher) Start(ctx context.Context) {
	slog.Info("Starting Job Dispatcher", slog.Int("workers", d.workers))

	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.workerLoop(ctx, i)
	}
}

func (d *defaultDispatcher) Stop() {
	slog.Info("Stopping Job Dispatcher...")
	close(d.stopChan)
	d.wg.Wait()
	slog.Info("Job Dispatcher stopped")
}

func (d *defaultDispatcher) workerLoop(ctx context.Context, workerID int) {
	defer d.wg.Done()

	for {
		select {
		case <-d.stopChan:
			return
		case <-ctx.Done():
			return
		default:
			// BRPOP blocks until a job is available or timeout (1 sec)
			res, err := d.redis.BRPop(ctx, 1*time.Second, defaultQueue).Result()
			if err != nil {
				if err != redis.Nil {
					// Unexpected error
					slog.Error("Failed to fetch job from redis", slog.String("error", err.Error()))
					time.Sleep(1 * time.Second)
				}
				continue
			}

			// res[0] is the key name, res[1] is the value
			if len(res) == 2 {
				d.processJob(ctx, res[1], workerID)
			}
		}
	}
}

func (d *defaultDispatcher) processJob(ctx context.Context, jobJSON string, workerID int) {
	var job model.Job
	if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
		slog.Error("Failed to unmarshal job", slog.String("error", err.Error()))
		return
	}

	handler, err := d.registry.Get(job.Type)
	if err != nil {
		slog.Error("No handler found for job type", slog.String("type", job.Type))
		d.failJob(ctx, &job, err.Error())
		return
	}

	slog.Info("Worker processing job", slog.Int("worker_id", workerID), slog.String("job_id", job.ID), slog.String("type", job.Type))

	// Mark as processing
	d.updateJobState(ctx, &job, model.JobStateProcessing, "")

	err = handler.Handle(ctx, &job)
	if err != nil {
		slog.Error("Job failed", slog.String("job_id", job.ID), slog.String("error", err.Error()))
		
		job.RetryCount++
		if job.RetryCount <= job.MaxRetries {
			d.updateJobState(ctx, &job, model.JobStateRetry, err.Error())
			// Push back to queue with exponential backoff logic (simplified for now to push immediately or to a delayed queue)
			// For this epic, we'll just push back to the tail of the main queue to be retried shortly.
			retryJSON, _ := json.Marshal(job)
			d.redis.RPush(ctx, defaultQueue, retryJSON)
		} else {
			d.failJob(ctx, &job, err.Error())
		}
		return
	}

	// Success
	slog.Info("Job completed successfully", slog.String("job_id", job.ID))
	d.updateJobState(ctx, &job, model.JobStateCompleted, "Success")
}

func (d *defaultDispatcher) updateJobState(ctx context.Context, job *model.Job, state model.JobState, msg string) {
	job.State = state
	
	updates := map[string]interface{}{
		"state":       state,
		"retry_count": job.RetryCount,
		"updated_at":  time.Now(),
	}
	if msg != "" {
		updates["error"] = msg
	}

	d.db.WithContext(ctx).Model(&model.Job{}).Where("id = ?", job.ID).Updates(updates)

	log := &model.JobLog{
		JobID:   uuid.MustParse(job.ID),
		State:   state,
		Message: msg,
	}
	d.db.WithContext(ctx).Create(log)
}

func (d *defaultDispatcher) failJob(ctx context.Context, job *model.Job, reason string) {
	d.updateJobState(ctx, job, model.JobStateFailed, reason)
}
