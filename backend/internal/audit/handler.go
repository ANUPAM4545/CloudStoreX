package audit

import (
	"net/http"
	"strconv"

	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListLogs(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		response.Error(c, http.StatusBadRequest, "invalid_request", "workspace_id is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, total, err := h.service.ListLogs(c.Request.Context(), workspaceID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"items": logs,
		"total": total,
	})
}
