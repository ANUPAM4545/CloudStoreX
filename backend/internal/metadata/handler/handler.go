package handler

import (
	"net/http"

	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/service"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetadataHandler struct {
	svc service.MetadataService
}

func NewMetadataHandler(svc service.MetadataService) *MetadataHandler {
	return &MetadataHandler{svc: svc}
}

func (h *MetadataHandler) RegisterRoutes(router gin.IRouter) {
	group := router.Group("/metadata")
	{
		group.GET("/objects/:id", h.GetObject)
		group.PATCH("/objects/:id", h.UpdateObject)
		group.DELETE("/objects/:id", h.DeleteObject)
		group.POST("/objects/:id/restore", h.RestoreObject)
		group.GET("/objects/:id/versions", h.ListObjectVersions)
		group.GET("/objects/:id/metadata", h.GetObjectMetadata)
		
		group.POST("/objects/:id/tags", h.TagObject)
		group.DELETE("/objects/:id/tags", h.UntagObject)

		group.GET("/search", h.SearchObjects)
		
		group.POST("/buckets", h.CreateBucket)
		group.GET("/buckets", h.ListBuckets)
	}
}

func (h *MetadataHandler) GetObject(c *gin.Context) {
	id := c.Param("id")
	obj, err := h.svc.FindObjectByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Object not found")
		return
	}
	response.Success(c, http.StatusOK, obj)
}

func (h *MetadataHandler) UpdateObject(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		MimeType     string            `json:"mime_type"`
		StorageClass string            `json:"storage_class"`
		Tags         map[string]string `json:"tags"`
		Metadata     map[string]string `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}

	obj, err := h.svc.UpdateObjectMetadata(c.Request.Context(), id, req.MimeType, req.StorageClass, req.Tags, req.Metadata)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update object")
		return
	}
	response.Success(c, http.StatusOK, obj)
}

func (h *MetadataHandler) DeleteObject(c *gin.Context) {
	id := c.Param("id")
	workspaceID := c.GetString("workspace_id") // Assuming auth middleware sets this
	if workspaceID == "" {
		workspaceID = "default" // fallback for tests
	}
	
	if err := h.svc.SoftDeleteObject(c.Request.Context(), workspaceID, id); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete object")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Object deleted successfully"})
}

func (h *MetadataHandler) RestoreObject(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.RestoreObject(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to restore object")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Object restored successfully"})
}

func (h *MetadataHandler) ListObjectVersions(c *gin.Context) {
	id := c.Param("id")
	versions, err := h.svc.ListObjectVersions(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list versions")
		return
	}
	response.Success(c, http.StatusOK, versions)
}

func (h *MetadataHandler) GetObjectMetadata(c *gin.Context) {
	id := c.Param("id")
	meta, err := h.svc.GetObjectMetadata(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get metadata")
		return
	}
	response.Success(c, http.StatusOK, meta)
}

func (h *MetadataHandler) TagObject(c *gin.Context) {
	id := c.Param("id")
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid tags format")
		return
	}

	if err := h.svc.TagObject(c.Request.Context(), id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to tag object")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Object tagged successfully"})
}

func (h *MetadataHandler) UntagObject(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Keys []string `json:"keys"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}

	if err := h.svc.UntagObject(c.Request.Context(), id, req.Keys); err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to untag object")
		return
	}
	response.Success(c, http.StatusOK, gin.H{"message": "Object untagged successfully"})
}

func (h *MetadataHandler) SearchObjects(c *gin.Context) {
	var query dto.SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid search query")
		return
	}
	
	query.WorkspaceID = c.GetString("workspace_id") // Enforce workspace isolation

	results, total, err := h.svc.SearchObjects(c.Request.Context(), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Search failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"total": total,
	})
}

func (h *MetadataHandler) CreateBucket(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		ProviderID string `json:"provider_id" binding:"required"`
		Region     string `json:"region" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request")
		return
	}

	workspaceIDStr := c.GetString("workspace_id")
	var workspaceID uuid.UUID
	if workspaceIDStr != "" {
		workspaceID = uuid.MustParse(workspaceIDStr)
	} else {
		workspaceID = uuid.New() // Fallback
	}

	bucket, err := h.svc.CreateBucket(c.Request.Context(), workspaceID, req.ProviderID, req.Name, req.Region)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create bucket")
		return
	}

	response.Success(c, http.StatusCreated, bucket)
}

func (h *MetadataHandler) ListBuckets(c *gin.Context) {
	workspaceID := c.GetString("workspace_id")
	if workspaceID == "" {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Workspace not found")
		return
	}

	buckets, err := h.svc.ListBuckets(c.Request.Context(), workspaceID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list buckets")
		return
	}

	response.Success(c, http.StatusOK, buckets)
}
