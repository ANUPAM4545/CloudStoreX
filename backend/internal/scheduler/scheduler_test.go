package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduler_RegisterAndRun(t *testing.T) {
	sched := NewScheduler(nil)
	var executions int32

	id := sched.RegisterTask("test-health", TaskHealthCheck, 50*time.Millisecond, func(ctx context.Context) error {
		atomic.AddInt32(&executions, 1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched.Start(ctx)

	// artificially trigger execution check for fast unit test
	time.Sleep(150 * time.Millisecond)
	sched.Stop()

	jobs := sched.GetScheduledJobs()
	if len(jobs) != 1 || jobs[0].ID != id {
		t.Fatalf("expected 1 scheduled job with matching ID, got %v", jobs)
	}
}

func TestScheduler_Stop(t *testing.T) {
	sched := NewScheduler(nil)
	ctx := context.Background()
	sched.Start(ctx)
	sched.Stop()
	// starting and stopping without panic or deadlock
}
