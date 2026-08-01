package workers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/audit"
	"github.com/cloudstorex/backend/internal/audit/model"
	jobmodel "github.com/cloudstorex/backend/internal/jobs/model"
	"github.com/cloudstorex/backend/internal/metadata/events"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditWorker consumes metadata and compliance events and persists them as AuditLog records.
type AuditWorker struct {
	auditService audit.Service
	log          *slog.Logger
}

func NewAuditWorker(auditService audit.Service, log *slog.Logger) *AuditWorker {
	if log == nil {
		log = slog.Default()
	}
	return &AuditWorker{
		auditService: auditService,
		log:          log,
	}
}

func (w *AuditWorker) Handle(ctx context.Context, job *jobmodel.Job) error {
	var event events.Event
	if err := json.Unmarshal([]byte(job.Payload), &event); err != nil {
		w.log.ErrorContext(ctx, "failed to unmarshal event for audit worker", slog.String("error", err.Error()))
		return err
	}

	workspaceID, _ := event.Payload["workspace_id"].(string)
	if workspaceID == "" {
		workspaceID = "default-workspace" // Fallback if workspace is not in event payload
	}

	detailsBytes, _ := json.Marshal(event.Payload)

	logEntry := &model.AuditLog{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Action:      string(event.Type),
		Resource:    event.ObjectID,
		Details:     datatypes.JSON(detailsBytes),
		Timestamp:   time.Now(),
	}

	if err := w.auditService.Record(ctx, logEntry); err != nil {
		return err
	}

	return nil
}
