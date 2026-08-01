package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// mockService implements storage.Service for handler testing.
type mockService struct {
	buckets map[string]time.Time
	objects map[string]*Object
	data    map[string][]byte
}

func newMockService() *mockService {
	return &mockService{
		buckets: make(map[string]time.Time),
		objects: make(map[string]*Object),
		data:    make(map[string][]byte),
	}
}

func (m *mockService) UploadObject(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	m.data[bucket+"/"+key] = content
	m.objects[bucket+"/"+key] = &Object{
		Key:          key,
		Bucket:       bucket,
		Size:         int64(len(content)),
		ETag:         "etag-" + key,
		LastModified: time.Now().UTC(),
		Metadata:     meta,
	}
	return &StorageResponse{
		Bucket: bucket,
		Key:    key,
		Size:   int64(len(content)),
		ETag:   "etag-" + key,
	}, nil
}

func (m *mockService) DownloadObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	content, ok := m.data[bucket+"/"+key]
	if !ok {
		return nil, ErrObjectNotFound
	}
	return io.NopCloser(bytes.NewReader(content)), nil
}

func (m *mockService) DeleteObject(ctx context.Context, bucket, key string) error {
	delete(m.data, bucket+"/"+key)
	delete(m.objects, bucket+"/"+key)
	return nil
}

func (m *mockService) ObjectExists(ctx context.Context, bucket, key string) (bool, error) {
	_, ok := m.objects[bucket+"/"+key]
	return ok, nil
}

func (m *mockService) ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error) {
	var list []*Object
	for _, o := range m.objects {
		if o.Bucket == bucket && strings.HasPrefix(o.Key, prefix) {
			list = append(list, o)
		}
	}
	return list, nil
}

func (m *mockService) CreateBucket(ctx context.Context, bucket string) error {
	m.buckets[bucket] = time.Now().UTC()
	return nil
}

func (m *mockService) DeleteBucket(ctx context.Context, bucket string) error {
	delete(m.buckets, bucket)
	return nil
}

func (m *mockService) ListBuckets(ctx context.Context) ([]*Bucket, error) {
	var list []*Bucket
	for name, created := range m.buckets {
		list = append(list, &Bucket{Name: name, CreatedAt: created})
	}
	return list, nil
}

func (m *mockService) GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "http://presigned-upload", nil
}

func (m *mockService) GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "http://presigned-download", nil
}

func (m *mockService) SearchObjects(ctx context.Context, query dto.SearchQuery) ([]dto.ObjectMetaDTO, int64, error) {
	return nil, 0, nil
}

func (m *mockService) GetObjectByID(ctx context.Context, id string) (*dto.ObjectMetaDTO, error) {
	return nil, nil
}

func (m *mockService) GetObjectMetadata(ctx context.Context, id string) (map[string]string, error) {
	return nil, nil
}

func (m *mockService) TagObject(ctx context.Context, id string, tags map[string]string) error {
	return nil
}

func (m *mockService) UntagObject(ctx context.Context, id string, keys []string) error {
	return nil
}

func (m *mockService) ListObjectVersions(ctx context.Context, id string) ([]dto.ObjectVersionDTO, error) {
	return nil, nil
}

func (m *mockService) RestoreObject(ctx context.Context, id string) error {
	return nil
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("request_id", "test-request-id-12345")
		c.Next()
	})
	api := r.Group("/api/v1/storage")
	{
		api.POST("/buckets", handler.CreateBucket)
		api.GET("/buckets", handler.ListBuckets)
		api.DELETE("/buckets/:bucket", handler.DeleteBucket)

		api.POST("/buckets/:bucket/objects", handler.UploadObject)
		api.GET("/buckets/:bucket/objects", handler.ListObjects)
		api.GET("/buckets/:bucket/objects/*key", handler.DownloadObject)
		api.DELETE("/buckets/:bucket/objects/*key", handler.DeleteObject)
		api.HEAD("/buckets/:bucket/objects/*key", handler.ObjectExists)
	}
	return r
}

func TestHandler_CreateBucket(t *testing.T) {
	svc := newMockService()
	h := NewHandler(svc, 100)
	r := setupTestRouter(h)

	// Valid bucket
	body := `{"bucket": "test-bucket-1"}`
	req := httptest.NewRequest("POST", "/api/v1/storage/buckets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
	if w.Header().Get("X-Request-ID") != "test-request-id-12345" {
		t.Errorf("expected X-Request-ID header to be set")
	}

	var resp response.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || !resp.Success {
		t.Fatalf("expected success response, got %s", w.Body.String())
	}

	// Invalid bucket name
	bodyInvalid := `{"bucket": "-invalid"}`
	reqInvalid := httptest.NewRequest("POST", "/api/v1/storage/buckets", strings.NewReader(bodyInvalid))
	reqInvalid.Header.Set("Content-Type", "application/json")
	wInvalid := httptest.NewRecorder()

	r.ServeHTTP(wInvalid, reqInvalid)

	if wInvalid.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid bucket name, got %d", wInvalid.Code)
	}
}

func TestHandler_ListAndDeleteBucket(t *testing.T) {
	svc := newMockService()
	h := NewHandler(svc, 100)
	r := setupTestRouter(h)

	svc.CreateBucket(context.Background(), "bucket-a")

	req := httptest.NewRequest("GET", "/api/v1/storage/buckets", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	reqDel := httptest.NewRequest("DELETE", "/api/v1/storage/buckets/bucket-a", nil)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)

	if wDel.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", wDel.Code)
	}
}

func TestHandler_UploadAndDownloadObject(t *testing.T) {
	svc := newMockService()
	h := NewHandler(svc, 100)
	r := setupTestRouter(h)

	// Create multipart stream
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("hello world storage API"))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/storage/buckets/mybucket/objects", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created on upload, got %d: %s", w.Code, w.Body.String())
	}

	// Download the object
	reqDown := httptest.NewRequest("GET", "/api/v1/storage/buckets/mybucket/objects/testfile.txt", nil)
	wDown := httptest.NewRecorder()

	r.ServeHTTP(wDown, reqDown)

	if wDown.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on download, got %d: %s", wDown.Code, wDown.Body.String())
	}
	if wDown.Body.String() != "hello world storage API" {
		t.Errorf("expected body %q, got %q", "hello world storage API", wDown.Body.String())
	}
	if wDown.Header().Get("ETag") == "" {
		t.Errorf("expected ETag header on download")
	}
	if wDown.Header().Get("Last-Modified") == "" {
		t.Errorf("expected Last-Modified header on download")
	}
	if wDown.Header().Get("Cache-Control") == "" {
		t.Errorf("expected Cache-Control header on download")
	}
}

func TestHandler_ObjectExistsAndDelete(t *testing.T) {
	svc := newMockService()
	h := NewHandler(svc, 100)
	r := setupTestRouter(h)

	svc.UploadObject(context.Background(), "mybucket", "docs/file.txt", strings.NewReader("content"), 7, nil)

	// Exists (HEAD)
	reqHead := httptest.NewRequest("HEAD", "/api/v1/storage/buckets/mybucket/objects/docs/file.txt", nil)
	wHead := httptest.NewRecorder()
	r.ServeHTTP(wHead, reqHead)
	if wHead.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for existing object, got %d", wHead.Code)
	}

	// Non-existing object (HEAD)
	reqMissing := httptest.NewRequest("HEAD", "/api/v1/storage/buckets/mybucket/objects/missing.txt", nil)
	wMissing := httptest.NewRecorder()
	r.ServeHTTP(wMissing, reqMissing)
	if wMissing.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found for missing object, got %d", wMissing.Code)
	}

	// Delete object
	reqDel := httptest.NewRequest("DELETE", "/api/v1/storage/buckets/mybucket/objects/docs/file.txt", nil)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)
	if wDel.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for delete object, got %d", wDel.Code)
	}
}
