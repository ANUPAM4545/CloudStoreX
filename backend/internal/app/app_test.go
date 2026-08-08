package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	
	"testing"

	"github.com/cloudstorex/backend/internal/provider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLivenessProbes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := &App{
		Router: gin.New(),
	}
	app.setupRoutes()

	endpoints := []string{"/healthz", "/livez", "/api/v1/live"}
	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ep, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var resp map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			data, ok := resp["data"].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "up", data["status"])
		})
	}
}

func TestReadinessProbe_FailureWhenUninitialized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	app := &App{
		Router:          gin.New(),
		StorageRegistry: provider.NewRegistry(), // empty registry
	}
	app.setupRoutes()

	endpoints := []string{"/readyz", "/api/v1/ready"}
	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, ep, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			app.Router.ServeHTTP(rr, req)

			// Since DB/Redis are nil or unreachable, readyz should return 503 Service Unavailable
			assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

			var resp map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			assert.Equal(t, false, resp["success"])
			errObj, ok := resp["error"].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "SERVICE_UNAVAILABLE", errObj["code"])
		})
	}
}
