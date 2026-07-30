package middleware

import (
	"net/http"
	"strings"

	"github.com/cloudstorex/backend/internal/identity"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService *identity.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "unauthorized", "Invalid or expired token")
			c.Abort()
			return
		}

		// Set the user ID in the context
		c.Set("user_id", userID.String())
		c.Next()
	}
}
