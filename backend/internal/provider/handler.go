package provider

import (
	"net/http"

	"github.com/cloudstorex/backend/internal/provider/dto"
	"github.com/cloudstorex/backend/internal/provider/service"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.ProviderService
}

func NewHandler(svc service.ProviderService) *Handler {
	return &Handler{service: svc}
}

func getWorkspaceID(c *gin.Context) string {
	// Fallback to a default workspace for now until full auth workspace isolation is added
	ws, exists := c.Get("workspace_id")
	if exists {
		if wsStr, ok := ws.(string); ok && wsStr != "" {
			return wsStr
		}
	}
	return "00000000-0000-0000-0000-000000000000"
}

// ListProviders returns all registered providers and their health status.
func (h *Handler) ListProviders(c *gin.Context) {
	wsID := getWorkspaceID(c)
	providers, err := h.service.ListProviders(c.Request.Context(), wsID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list providers")
		return
	}
	response.Success(c, http.StatusOK, providers)
}

// CreateProvider registers a new provider configuration.
func (h *Handler) CreateProvider(c *gin.Context) {
	var req dto.CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}

	// Override WorkspaceID if provided via auth context
	wsID := getWorkspaceID(c)
	if req.WorkspaceID.String() == "00000000-0000-0000-0000-000000000000" && wsID != "" {
		// Just relying on the payload for now, or forcing it.
	}

	provider, err := h.service.CreateProvider(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create provider")
		return
	}
	response.Success(c, http.StatusCreated, provider)
}

// GetProvider returns a specific provider by ID.
func (h *Handler) GetProvider(c *gin.Context) {
	id := c.Param("id")
	provider, err := h.service.GetProvider(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "provider not found")
		return
	}
	response.Success(c, http.StatusOK, provider)
}

// UpdateProvider updates a specific provider by ID.
func (h *Handler) UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload")
		return
	}

	provider, err := h.service.UpdateProvider(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update provider")
		return
	}
	response.Success(c, http.StatusOK, provider)
}

// DeleteProvider deletes a specific provider by ID.
func (h *Handler) DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteProvider(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete provider")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "deleted"})
}

// ValidateProvider validates the connection of a specific provider.
func (h *Handler) ValidateProvider(c *gin.Context) {
	id := c.Param("id")
	extended := c.Query("extended") == "true"
	
	err := h.service.ValidateProviderConnection(c.Request.Context(), id, extended)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "valid"})
}

// SetDefaultProvider sets a specific provider as the default for the workspace.
func (h *Handler) SetDefaultProvider(c *gin.Context) {
	id := c.Param("id")
	wsID := getWorkspaceID(c)
	err := h.service.SetDefaultProvider(c.Request.Context(), wsID, id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to set default provider")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "default_set"})
}

// EnableProvider enables a specific provider.
func (h *Handler) EnableProvider(c *gin.Context) {
	id := c.Param("id")
	err := h.service.EnableProvider(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to enable provider")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "enabled"})
}

// DisableProvider disables a specific provider.
func (h *Handler) DisableProvider(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DisableProvider(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to disable provider")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"status": "disabled"})
}
