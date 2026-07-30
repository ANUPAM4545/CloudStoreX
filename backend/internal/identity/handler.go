package identity

import (
	"errors"
	"net/http"

	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type registerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			response.Error(c, http.StatusConflict, "user_exists", "A user with this email already exists.")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password.")
			return
		}
		response.Error(c, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"token": token,
	})
}
