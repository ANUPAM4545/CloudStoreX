package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudstorex/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.PrometheusMiddleware())
	r.GET("/api/v1/metrics-test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.GET("/api/v1/error-test", func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "Error")
	})

	// Test 200 OK request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics-test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test 500 Internal Server Error request
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/error-test", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)
}
