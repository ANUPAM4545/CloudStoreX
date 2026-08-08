package handler

import (
	"net/http"
	"strconv"

	"github.com/cloudstorex/backend/internal/policy/dto"
	"github.com/cloudstorex/backend/internal/policy/engine"
	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/cloudstorex/backend/internal/policy/service"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service service.PolicyService
	engine  engine.PolicyEngine
}

func NewHandler(svc service.PolicyService, eng engine.PolicyEngine) *Handler {
	return &Handler{service: svc, engine: eng}
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

type EvaluateRequest struct {
	Bucket       string            `json:"bucket"`
	ObjectKey    string            `json:"object_key"`
	Size         int64             `json:"size"`
	MimeType     string            `json:"mime_type"`
	StorageClass string            `json:"storage_class"`
	Tags         map[string]string `json:"tags"`
	Operation    string            `json:"operation"`
}

func (h *Handler) Evaluate(c *gin.Context) {
	var req EvaluateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}
	
	wsID := getWorkspaceID(c)
	wid, err := uuid.Parse(wsID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_WORKSPACE", "invalid workspace id")
		return
	}
	
	evalCtx := &model.EvaluationContext{
		WorkspaceID:  wid,
		Bucket:       req.Bucket,
		ObjectKey:    req.ObjectKey,
		Size:         req.Size,
		MimeType:     req.MimeType,
		StorageClass: req.StorageClass,
		Tags:         req.Tags,
	}
	
	_, decision, err := h.engine.Resolve(c.Request.Context(), evalCtx, req.Operation)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "evaluation failed")
		return
	}
	
	response.Success(c, http.StatusOK, decision)
}
