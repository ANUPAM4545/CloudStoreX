package analytics

import (
	"net/http"
	"time"

	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAnalyticsHistory(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		response.Error(c, http.StatusBadRequest, "invalid_request", "workspace_id is required")
		return
	}

	fromStr := c.DefaultQuery("from", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	toStr := c.DefaultQuery("to", time.Now().Format("2006-01-02"))

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid from date format")
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid to date format")
		return
	}

	snapshots, err := h.service.GetHistory(c.Request.Context(), workspaceID, from, to)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, snapshots)
}

func (h *Handler) GetLatestMetrics(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		response.Error(c, http.StatusBadRequest, "invalid_request", "workspace_id is required")
		return
	}

	snapshot, err := h.service.GetLatestMetrics(c.Request.Context(), workspaceID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "not_found", "metrics not found")
		return
	}

	response.Success(c, http.StatusOK, snapshot)
}
