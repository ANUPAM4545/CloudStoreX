package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudstorex/backend/internal/metadata/handler"
	"github.com/cloudstorex/backend/internal/metadata/model"
	"github.com/cloudstorex/backend/internal/metadata/service"
	"github.com/cloudstorex/backend/internal/metadata/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// A simple test using the existing mock to verify routing works.
type mockRepo struct {
	repository.MetadataRepository
}

func (m *mockRepo) FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error) {
	return &model.Bucket{
		ID:          uuid.New(),
		WorkspaceID: uuid.MustParse(workspaceID),
		Name:        bucketName,
	}, nil
}

func TestMetadataHandler_ListBuckets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Since we're just testing the handler, we can pass nil or simple mocks.
	// We'll skip deep mocking here for brevity, testing a failure path:
	svc := service.NewMetadataService(nil, nil, nil)
	h := handler.NewMetadataHandler(svc)

	router := gin.Default()
	// Mock auth middleware setting workspace_id
	router.Use(func(c *gin.Context) {
		c.Set("workspace_id", "") // Empty workspace ID should trigger 401
		c.Next()
	})
	h.RegisterRoutes(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metadata/buckets", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
