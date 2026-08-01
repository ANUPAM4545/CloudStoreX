package quota

import (
	"net/http"

	"github.com/cloudstorex/backend/internal/quota/model"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type SetQuotaRequest struct {
	MaxBytes   int64 `json:"max_bytes"`
	MaxObjects int64 `json:"max_objects"`
}

func (h *Handler) SetWorkspaceQuota(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		response.Error(c, http.StatusBadRequest, "invalid_request", "workspace_id is required")
		return
	}

	var req SetQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	quota := &model.WorkspaceQuota{
		WorkspaceID: workspaceID,
		MaxBytes:    req.MaxBytes,
		MaxObjects:  req.MaxObjects,
	}

	if err := h.service.SetWorkspaceQuota(c.Request.Context(), quota); err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, quota)
}

func (h *Handler) GetWorkspaceQuota(c *gin.Context) {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		response.Error(c, http.StatusBadRequest, "invalid_request", "workspace_id is required")
		return
	}

	quota, err := h.service.GetWorkspaceQuota(c.Request.Context(), workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, quota)
}
