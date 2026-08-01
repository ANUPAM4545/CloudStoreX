package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/analytics"
	"github.com/cloudstorex/backend/internal/analytics/model"
	jobmodel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/cloudstorex/backend/internal/quota"
	"github.com/google/uuid"
)

// AnalyticsWorker computes snapshots of storage usage and records them.
type AnalyticsWorker struct {
	analyticsService analytics.Service
	quotaService     quota.Service
	log              *slog.Logger
}

func NewAnalyticsWorker(analyticsService analytics.Service, quotaService quota.Service, log *slog.Logger) *AnalyticsWorker {
	if log == nil {
		log = slog.Default()
	}
	return &AnalyticsWorker{
		analyticsService: analyticsService,
		quotaService:     quotaService,
		log:              log,
	}
}

func (w *AnalyticsWorker) Handle(ctx context.Context, job *jobmodel.Job) error {
	w.log.InfoContext(ctx, "AnalyticsWorker started aggregation")

	// For simple multi-workspace analytics, we can compute snapshot from quota record
	// Here we'll use a placeholder "default-workspace" or query all quotas
	workspaceID := "default-workspace"
	q, err := w.quotaService.GetWorkspaceQuota(ctx, workspaceID)
	var totalBytes, totalObjects int64
	if err == nil && q != nil {
		totalBytes = q.BytesUsed
		totalObjects = q.ObjectCount
	}

	snapshot := &model.StorageSnapshot{
		ID:           uuid.New(),
		WorkspaceID:  workspaceID,
		TotalBytes:   totalBytes,
		TotalObjects: totalObjects,
		BucketCount:  1, // Default or computed from metadata service
		Date:         time.Now().Truncate(24 * time.Hour),
	}

	if err := w.analyticsService.RecordSnapshot(ctx, snapshot); err != nil {
		w.log.ErrorContext(ctx, "failed to record analytics snapshot", slog.String("error", err.Error()))
		return err
	}

	w.log.InfoContext(ctx, "AnalyticsWorker completed aggregation")
	return nil
}
