package handler

import (
	"net/http"
	"strconv"

	"github.com/cloudstorex/backend/internal/policy/dto"
	"github.com/cloudstorex/backend/internal/policy/service"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.PolicyService
}

func NewHandler(svc service.PolicyService) *Handler {
	return &Handler{service: svc}
}

func getWorkspaceID(c *gin.Context) string {
	ws, exists := c.Get("workspace_id")
	if exists {
		if wsStr, ok := ws.(string); ok && wsStr != "" {
			return wsStr
		}
	}
	return "00000000-0000-0000-0000-000000000000"
}

func (h *Handler) ListPolicies(c *gin.Context) {
	wsID := getWorkspaceID(c)
	policies, err := h.service.ListPolicies(c.Request.Context(), wsID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list policies")
		return
	}
	response.Success(c, http.StatusOK, policies)
}

func (h *Handler) CreatePolicy(c *gin.Context) {
	var req dto.CreatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}
	
	// Override WorkspaceID if provided via auth context
	wsID := getWorkspaceID(c)
	if req.WorkspaceID.String() == "00000000-0000-0000-0000-000000000000" && wsID != "" {
		// parse it
	}

	policy, err := h.service.CreatePolicy(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create policy")
		return
	}
	response.Success(c, http.StatusCreated, policy)
}

func (h *Handler) GetPolicy(c *gin.Context) {
	id := c.Param("id")
	policy, err := h.service.GetPolicy(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "policy not found")
		return
	}
	response.Success(c, http.StatusOK, policy)
}

func (h *Handler) UpdatePolicy(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}

	policy, err := h.service.UpdatePolicy(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update policy")
		return
	}
	response.Success(c, http.StatusOK, policy)
}

func (h *Handler) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeletePolicy(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete policy")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) EnablePolicy(c *gin.Context) {
	id := c.Param("id")
	err := h.service.EnablePolicy(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to enable policy")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "enabled"})
}

func (h *Handler) DisablePolicy(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DisablePolicy(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to disable policy")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "disabled"})
}

func (h *Handler) ListRoutingDecisions(c *gin.Context) {
	wsID := getWorkspaceID(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	decisions, total, err := h.service.ListRoutingDecisions(c.Request.Context(), wsID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list routing decisions")
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"items": decisions,
		"total": total,
	})
}
