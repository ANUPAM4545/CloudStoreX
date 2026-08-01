package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPprof_DisabledByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPprof(router, false, "")

	req, _ := http.NewRequest(http.MethodGet, "/debug/pprof/heap", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 Not Found when pprof disabled, got %d", w.Code)
	}
}

func TestPprof_EnabledWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPprof(router, true, "")

	req, _ := http.NewRequest(http.MethodGet, "/debug/pprof/cmdline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK when pprof enabled without token, got %d", w.Code)
	}
}

func TestPprof_EnabledWithAdminTokenSecurity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPprof(router, true, "secret-admin-token")

	// 1. Request without token -> 401
	reqUnauthorized, _ := http.NewRequest(http.MethodGet, "/debug/pprof/cmdline", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, reqUnauthorized)
	if w1.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized without token, got %d", w1.Code)
	}

	// 2. Request with invalid token -> 401
	reqInvalid, _ := http.NewRequest(http.MethodGet, "/debug/pprof/cmdline?token=wrong", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, reqInvalid)
	if w2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized with wrong token, got %d", w2.Code)
	}

	// 3. Request with valid header token -> 200
	reqValid, _ := http.NewRequest(http.MethodGet, "/debug/pprof/cmdline", nil)
	reqValid.Header.Set("X-Pprof-Token", "secret-admin-token")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, reqValid)
	if w3.Code != http.StatusOK {
		t.Fatalf("expected 200 OK with valid X-Pprof-Token header, got %d", w3.Code)
	}
}
