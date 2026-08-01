package lifecycle

import (
	"net/http"

	"github.com/cloudstorex/backend/internal/lifecycle/model"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRule(c *gin.Context) {
	bucketIDStr := c.Param("bucket_id")
	bucketID, err := uuid.Parse(bucketIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_bucket_id", "invalid bucket_id")
		return
	}

	var rule model.LifecycleRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}
	rule.BucketID = bucketID

	if err := h.service.CreateRule(c.Request.Context(), &rule); err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, rule)
}

func (h *Handler) ListRules(c *gin.Context) {
	bucketIDStr := c.Param("bucket_id")
	bucketID, err := uuid.Parse(bucketIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_bucket_id", "invalid bucket_id")
		return
	}

	rules, err := h.service.ListRules(c.Request.Context(), bucketID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	response.Success(c, http.StatusOK, rules)
}

func (h *Handler) DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_rule_id", "invalid rule id")
		return
	}

	if err := h.service.DeleteRule(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
