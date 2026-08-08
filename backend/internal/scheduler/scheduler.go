package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskHandler func(ctx context.Context) error

// Scheduler coordinates recurring reliability jobs across workers (Refinement 4).
type Scheduler interface {
	RegisterTask(name string, taskType TaskType, interval time.Duration, handler TaskHandler) uuid.UUID
	Start(ctx context.Context)
	Stop()
	GetScheduledJobs() []ScheduledJob
}

type taskEntry struct {
	job     ScheduledJob
	handler TaskHandler
}

type defaultScheduler struct {
	mu      sync.RWMutex
	tasks   map[uuid.UUID]*taskEntry
	logger  *slog.Logger
	stopCh  chan struct{}
	running bool
}

// NewScheduler creates a new thread-safe Central Scheduler.
func NewScheduler(logger *slog.Logger) Scheduler {
	if logger == nil {
		logger = slog.Default()
	}
	return &defaultScheduler{
		tasks:  make(map[uuid.UUID]*taskEntry),
		logger: logger,
		stopCh: make(chan struct{}),
	}
}

func (s *defaultScheduler) RegisterTask(name string, taskType TaskType, interval time.Duration, handler TaskHandler) uuid.UUID {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New()
	now := time.Now().UTC()
	s.tasks[id] = &taskEntry{
		job: ScheduledJob{
			ID:        id,
			Name:      name,
			TaskType:  taskType,
			Interval:  interval,
			LastRunAt: time.Time{},
			NextRunAt: now.Add(interval),
			Status:    "IDLE",
			RunCount:  0,
		},
		handler: handler,
	}
	s.logger.Info("registered central scheduler task",
		slog.String("id", id.String()),
		slog.String("name", name),
		slog.String("task_type", string(taskType)),
		slog.Duration("interval", interval),
	)
	return id
}

func (s *defaultScheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stopCh = make(chan struct{})
	s.mu.Unlock()

	go s.runLoop(ctx)
}

func (s *defaultScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	close(s.stopCh)
	s.running = false
}

func (s *defaultScheduler) GetScheduledJobs() []ScheduledJob {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]ScheduledJob, 0, len(s.tasks))
	for _, entry := range s.tasks {
		jobs = append(jobs, entry.job)
	}
	return jobs
}

func (s *defaultScheduler) runLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.checkAndRun(ctx)
		}
	}
}

func (s *defaultScheduler) checkAndRun(ctx context.Context) {
	now := time.Now().UTC()

	s.mu.Lock()
	var toRun []*taskEntry
	for _, entry := range s.tasks {
		if now.After(entry.job.NextRunAt) || now.Equal(entry.job.NextRunAt) {
			entry.job.Status = "RUNNING"
			toRun = append(toRun, entry)
		}
	}
	s.mu.Unlock()

	for _, entry := range toRun {
		go func(e *taskEntry) {
			err := e.handler(ctx)
			s.mu.Lock()
			defer s.mu.Unlock()
			e.job.LastRunAt = time.Now().UTC()
			e.job.NextRunAt = time.Now().UTC().Add(e.job.Interval)
			e.job.Status = "IDLE"
			e.job.RunCount++
			if err != nil {
				s.logger.Error("central scheduler task execution failed",
					slog.String("task_name", e.job.Name),
					slog.String("error", err.Error()),
				)
			}
		}(entry)
	}
}
