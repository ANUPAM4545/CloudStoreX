package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudstorex/backend/internal/chaos/model"
	"github.com/cloudstorex/backend/internal/chaos/service"
	"github.com/google/uuid"
)

type dummyChaosService struct {
	lastExp *model.ChaosExperiment
}

func (d *dummyChaosService) RunExperiment(ctx context.Context, req service.RunExperimentRequest) (*model.ChaosExperiment, error) {
	exp := &model.ChaosExperiment{
		ID:             uuid.New(),
		WorkspaceID:    req.WorkspaceID,
		Name:           req.Name,
		ExperimentType: req.ExperimentType,
		Status:         model.StatusRunning,
	}
	d.lastExp = exp
	return exp, nil
}
func (d *dummyChaosService) CompleteExperiment(ctx context.Context, expID uuid.UUID, measuredRPO, measuredRTO int) error {
	return nil
}
func (d *dummyChaosService) AbortExperiment(ctx context.Context, expID uuid.UUID, reason string) error {
	return nil
}
func (d *dummyChaosService) GetExperiment(ctx context.Context, expID uuid.UUID) (*model.ChaosExperiment, error) {
	return nil, nil
}
func (d *dummyChaosService) ListExperiments(ctx context.Context, workspaceID uuid.UUID) ([]model.ChaosExperiment, error) {
	return nil, nil
}

func TestHandler_RunExperiment(t *testing.T) {
	svc := &dummyChaosService{}
	h := NewHandler(svc)

	reqBody := service.RunExperimentRequest{
		WorkspaceID:     uuid.New(),
		Name:            "Test API Outage",
		ExperimentType:  model.TypeProviderOutage,
		TargetProvider:  "aws-s3-east",
		DurationSeconds: 60,
	}
	data, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/chaos/run", bytes.NewReader(data))
	rec := httptest.NewRecorder()

	h.RunExperiment(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", rec.Code)
	}
	if svc.lastExp == nil || svc.lastExp.Name != "Test API Outage" {
		t.Fatalf("expected experiment to be started via HTTP API")
	}
}
