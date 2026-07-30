package app

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/storage"
)

func getAuthToken(t *testing.T, a *App) string {
	t.Helper()

	email := fmtSprintf("test_storage_%d@example.com", time.Now().UnixNano())
	body := `{"email":"` + email + `","password":"password123","full_name":"Storage Test User"}`

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		// Try login if already exists
		loginBody := `{"email":"` + email + `","password":"password123"}`
		reqLogin := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(loginBody))
		reqLogin.Header.Set("Content-Type", "application/json")
		wLogin := httptest.NewRecorder()
		a.Router.ServeHTTP(wLogin, reqLogin)
		if wLogin.Code != http.StatusOK {
			t.Fatalf("failed to authenticate test user: %s", wLogin.Body.String())
		}
		var authResp struct {
			Data struct {
				Token string `json:"token"`
			} `json:"data"`
		}
		_ = json.Unmarshal(wLogin.Body.Bytes(), &authResp)
		return authResp.Data.Token
	}

	var authResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &authResp); err != nil {
		t.Fatalf("failed to parse auth response: %v", err)
	}
	return authResp.Data.Token
}

func fmtSprintf(format string, args ...interface{}) string {
	// Simple string formatter helper for test emails
	return strings.Replace(format, "%d", strings.Replace(time.Now().Format("150405.000"), ".", "", 1), 1)
}

func TestStorageAPI_EndToEndIntegration(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Skipf("skipping TestStorageAPI_EndToEndIntegration: database or MinIO unavailable: %v", err)
		return
	}

	token := getAuthToken(t, app)
	if token == "" {
		t.Skip("skipping test: could not obtain JWT token")
		return
	}
	authHeader := "Bearer " + token

	bucketName := "cloudstorex-api-e2e"

	// 1. Create Bucket
	reqCreateBucket := httptest.NewRequest("POST", "/api/v1/storage/buckets", strings.NewReader(`{"bucket":"`+bucketName+`"}`))
	reqCreateBucket.Header.Set("Content-Type", "application/json")
	reqCreateBucket.Header.Set("Authorization", authHeader)
	wCreateBucket := httptest.NewRecorder()

	app.Router.ServeHTTP(wCreateBucket, reqCreateBucket)
	if wCreateBucket.Code != http.StatusCreated && wCreateBucket.Code != http.StatusConflict {
		t.Fatalf("expected 201 Created or 409 Conflict on CreateBucket, got %d: %s", wCreateBucket.Code, wCreateBucket.Body.String())
	}

	// 2. List Buckets
	reqListBuckets := httptest.NewRequest("GET", "/api/v1/storage/buckets", nil)
	reqListBuckets.Header.Set("Authorization", authHeader)
	wListBuckets := httptest.NewRecorder()
	app.Router.ServeHTTP(wListBuckets, reqListBuckets)

	if wListBuckets.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on ListBuckets, got %d: %s", wListBuckets.Code, wListBuckets.Body.String())
	}
	var listBucketsResp struct {
		Success bool                 `json:"success"`
		Data    []*storage.BucketDTO `json:"data"`
	}
	if err := json.Unmarshal(wListBuckets.Body.Bytes(), &listBucketsResp); err != nil || !listBucketsResp.Success {
		t.Fatalf("expected valid list buckets JSON response, got %s", wListBuckets.Body.String())
	}

	// 3. Upload Object via Multipart Form-Data Stream
	var buf bytes.Buffer
	mpWriter := multipart.NewWriter(&buf)
	part, err := mpWriter.CreateFormFile("file", "api-e2e-test.txt")
	if err != nil {
		t.Fatalf("failed to create multipart file: %v", err)
	}
	_, _ = part.Write([]byte("CloudStoreX Epic 4 End-to-End REST API integration test!"))
	_ = mpWriter.Close()

	reqUpload := httptest.NewRequest("POST", "/api/v1/storage/buckets/"+bucketName+"/objects", &buf)
	reqUpload.Header.Set("Content-Type", mpWriter.FormDataContentType())
	reqUpload.Header.Set("Authorization", authHeader)
	wUpload := httptest.NewRecorder()

	app.Router.ServeHTTP(wUpload, reqUpload)
	if wUpload.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on UploadObject, got %d: %s", wUpload.Code, wUpload.Body.String())
	}
	if wUpload.Header().Get("X-Request-ID") == "" {
		t.Errorf("expected X-Request-ID header on UploadObject response")
	}

	// 4. Object Exists (HEAD)
	reqHead := httptest.NewRequest("HEAD", "/api/v1/storage/buckets/"+bucketName+"/objects/api-e2e-test.txt", nil)
	reqHead.Header.Set("Authorization", authHeader)
	wHead := httptest.NewRecorder()

	app.Router.ServeHTTP(wHead, reqHead)
	if wHead.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on HEAD ObjectExists, got %d", wHead.Code)
	}

	// 5. Download Object (GET stream)
	reqDownload := httptest.NewRequest("GET", "/api/v1/storage/buckets/"+bucketName+"/objects/api-e2e-test.txt", nil)
	reqDownload.Header.Set("Authorization", authHeader)
	wDownload := httptest.NewRecorder()

	app.Router.ServeHTTP(wDownload, reqDownload)
	if wDownload.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on DownloadObject, got %d: %s", wDownload.Code, wDownload.Body.String())
	}
	if wDownload.Body.String() != "CloudStoreX Epic 4 End-to-End REST API integration test!" {
		t.Errorf("unexpected download content: %s", wDownload.Body.String())
	}
	if wDownload.Header().Get("ETag") == "" {
		t.Errorf("expected ETag header on download response")
	}
	if wDownload.Header().Get("Cache-Control") == "" {
		t.Errorf("expected Cache-Control header on download response")
	}
	if !strings.Contains(wDownload.Header().Get("Content-Disposition"), "api-e2e-test.txt") {
		t.Errorf("expected Content-Disposition header with filename, got %s", wDownload.Header().Get("Content-Disposition"))
	}

	// 6. List Objects (GET)
	reqListObjects := httptest.NewRequest("GET", "/api/v1/storage/buckets/"+bucketName+"/objects?prefix=api-e2e", nil)
	reqListObjects.Header.Set("Authorization", authHeader)
	wListObjects := httptest.NewRecorder()

	app.Router.ServeHTTP(wListObjects, reqListObjects)
	if wListObjects.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on ListObjects, got %d: %s", wListObjects.Code, wListObjects.Body.String())
	}
	var listObjectsResp struct {
		Success bool                 `json:"success"`
		Data    []*storage.ObjectDTO `json:"data"`
	}
	if err := json.Unmarshal(wListObjects.Body.Bytes(), &listObjectsResp); err != nil || !listObjectsResp.Success {
		t.Fatalf("expected valid list objects JSON response, got %s", wListObjects.Body.String())
	}

	// 7. Delete Object (DELETE)
	reqDeleteObject := httptest.NewRequest("DELETE", "/api/v1/storage/buckets/"+bucketName+"/objects/api-e2e-test.txt", nil)
	reqDeleteObject.Header.Set("Authorization", authHeader)
	wDeleteObject := httptest.NewRecorder()

	app.Router.ServeHTTP(wDeleteObject, reqDeleteObject)
	if wDeleteObject.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on DeleteObject, got %d: %s", wDeleteObject.Code, wDeleteObject.Body.String())
	}

	// 8. Verify object is deleted
	reqHeadAfterDel := httptest.NewRequest("HEAD", "/api/v1/storage/buckets/"+bucketName+"/objects/api-e2e-test.txt", nil)
	reqHeadAfterDel.Header.Set("Authorization", authHeader)
	wHeadAfterDel := httptest.NewRecorder()

	app.Router.ServeHTTP(wHeadAfterDel, reqHeadAfterDel)
	if wHeadAfterDel.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found after deleting object, got %d", wHeadAfterDel.Code)
	}

	// 9. Delete Bucket (DELETE)
	reqDeleteBucket := httptest.NewRequest("DELETE", "/api/v1/storage/buckets/"+bucketName, nil)
	reqDeleteBucket.Header.Set("Authorization", authHeader)
	wDeleteBucket := httptest.NewRecorder()

	app.Router.ServeHTTP(wDeleteBucket, reqDeleteBucket)
	if wDeleteBucket.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content on DeleteBucket, got %d: %s", wDeleteBucket.Code, wDeleteBucket.Body.String())
	}
}
