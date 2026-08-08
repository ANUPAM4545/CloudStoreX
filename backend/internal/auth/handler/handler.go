package handler

import (
	"net/http"

	"github.com/cloudstorex/backend/internal/auth/model"
	"github.com/cloudstorex/backend/internal/auth/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService service.Service
}

func NewHandler(authService service.Service) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.GET("/login/:provider", h.Login)
		authGroup.GET("/callback/:provider", h.Callback)
	}
}

func (h *Handler) Login(c *gin.Context) {
	providerStr := c.Param("provider")
	ptype := model.ProviderType(providerStr)

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:3000/dashboard"
	}

	url, err := h.authService.GetLoginURL(c.Request.Context(), ptype, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) Callback(c *gin.Context) {
	providerStr := c.Param("provider")
	ptype := model.ProviderType(providerStr)

	res, err := h.authService.HandleCallback(c.Request.Context(), ptype, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	}

	// For now, return the mock result.
	// A real application would issue a Session Token (or JWT) to the client here via Set-Cookie.
	c.JSON(http.StatusOK, res)
}
