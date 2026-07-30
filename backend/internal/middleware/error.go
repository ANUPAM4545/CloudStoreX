package middleware

import (
	"net/http"

	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// We could log the stack trace here
				response.Error(c, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
				c.Abort()
			}
		}()
		c.Next()
	}
}
