package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpsHandler_GetClusterHealth(t *testing.T) {
	h := NewOpsHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/operations/cluster/health", nil)
	rec := httptest.NewRecorder()

	h.GetClusterHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}
