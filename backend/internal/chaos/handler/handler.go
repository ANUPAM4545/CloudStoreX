package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cloudstorex/backend/internal/chaos/service"
	"github.com/google/uuid"
)

// Handler provides HTTP endpoints for managing and inspecting chaos experiments.
type Handler struct {
	svc service.Service
}

// NewHandler creates a new Phase 8 Chaos Engineering HTTP handler.
func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// RunExperiment handles POST requests to start a controlled chaos experiment.
func (h *Handler) RunExperiment(w http.ResponseWriter, r *http.Request) {
	var req service.RunExperimentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exp, err := h.svc.RunExperiment(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(exp)
}

type completeReq struct {
	MeasuredRPO int `json:"measured_rpo"`
	MeasuredRTO int `json:"measured_rto"`
}

// CompleteExperiment handles POST requests to complete an experiment and evaluate RTO/RPO SLA adherence.
func (h *Handler) CompleteExperiment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	expID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}

	var req completeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.svc.CompleteExperiment(r.Context(), expID, req.MeasuredRPO, req.MeasuredRTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
