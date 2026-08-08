package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskReplicationVerify TaskType = "REPLICATION_VERIFY"
	TaskBackupSchedule    TaskType = "BACKUP_SCHEDULE"
	TaskConsistencyCheck  TaskType = "CONSISTENCY_CHECK"
	TaskRepairSweep       TaskType = "REPAIR_SWEEP"
	TaskHealthCheck       TaskType = "HEALTH_CHECK"
	TaskLifecycleSweep    TaskType = "LIFECYCLE_SWEEP"
)

// ScheduledJob represents a repeating reliability task managed by the Central Scheduler.
type ScheduledJob struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	TaskType  TaskType      `json:"task_type"`
	Interval  time.Duration `json:"interval"`
	LastRunAt time.Time     `json:"last_run_at"`
	NextRunAt time.Time     `json:"next_run_at"`
	Status    string        `json:"status"`
	RunCount  int64         `json:"run_count"`
}
