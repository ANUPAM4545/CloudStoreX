package middleware

import (
	"net/http"
	"strings"

	"github.com/cloudstorex/backend/internal/identity"
	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func AuthMiddleware(tokenService *identity.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartChildSpan(c.Request.Context(), "AuthMiddleware.ValidateToken")
		defer span.End()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			tracing.RecordError(span, http.ErrNoCookie)
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
			tracing.RecordError(span, err)
			response.Error(c, http.StatusUnauthorized, "unauthorized", "Invalid or expired token")
			c.Abort()
			return
		}

		span.SetAttributes(attribute.String("auth.user_id", userID.String()))
		c.Request = c.Request.WithContext(ctx)

		// Set the user ID in the context
		c.Set("user_id", userID.String())
		c.Next()
	}
}
