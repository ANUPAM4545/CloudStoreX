package storage

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/cloudstorex/backend/internal/storage/validator"
	"github.com/gin-gonic/gin"
)

// Handler exposes authenticated REST endpoints for bucket and object storage operations.
type Handler struct {
	service       Service
	maxUploadSize int64 // Maximum upload size in bytes
}

// NewHandler creates a new Storage HTTP Handler.
func NewHandler(service Service, maxUploadSizeMB int64) *Handler {
	maxBytes := maxUploadSizeMB * 1024 * 1024
	if maxBytes <= 0 {
		maxBytes = 100 * 1024 * 1024 // 100MB default
	}
	return &Handler{
		service:       service,
		maxUploadSize: maxBytes,
	}
}

func getRequestID(c *gin.Context) string {
	reqID := c.Writer.Header().Get("X-Request-ID")
	if reqID == "" {
		reqID = c.GetString("request_id")
	}
	return reqID
}

func setResponseHeaders(c *gin.Context) {
	reqID := getRequestID(c)
	if reqID != "" && c.Writer.Header().Get("X-Request-ID") == "" {
		c.Header("X-Request-ID", reqID)
	}
}

func (h *Handler) handleStorageError(c *gin.Context, err error) {
	setResponseHeaders(c)

	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		err = domainErr.Err
	}

	switch {
	case errors.Is(err, ErrBucketExists):
		response.Error(c, http.StatusConflict, "bucket_exists", "A bucket with this name already exists.")
	case errors.Is(err, ErrBucketNotFound):
		response.Error(c, http.StatusNotFound, "bucket_not_found", "The specified bucket does not exist.")
	case errors.Is(err, ErrObjectNotFound):
		response.Error(c, http.StatusNotFound, "object_not_found", "The specified object does not exist.")
	case errors.Is(err, ErrInvalidRequest):
		response.Error(c, http.StatusBadRequest, "invalid_request", "The request parameters are invalid.")
	case errors.Is(err, ErrUnsupportedFeature):
		response.Error(c, http.StatusNotImplemented, "unsupported_feature", "The requested operation is not supported.")
	case errors.Is(err, ErrUploadFailed):
		response.Error(c, http.StatusInternalServerError, "upload_failed", "Failed to upload object.")
	case errors.Is(err, ErrDownloadFailed):
		response.Error(c, http.StatusInternalServerError, "download_failed", "Failed to download object.")
	case errors.Is(err, ErrDeleteFailed):
		response.Error(c, http.StatusInternalServerError, "delete_failed", "Failed to delete resource.")
	default:
		// Check validator message
		errStr := err.Error()
		if strings.Contains(errStr, "bucket name") || strings.Contains(errStr, "object key") || strings.Contains(errStr, "file size") {
			response.Error(c, http.StatusBadRequest, "validation_error", errStr)
			return
		}
		response.Error(c, http.StatusInternalServerError, "storage_error", "An unexpected storage error occurred.")
	}
}

// CreateBucket godoc
// @Summary Create a storage bucket
// @Description Creates a new bucket in the configured default storage provider
// @Tags buckets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBucketRequest true "Bucket creation request"
// @Success 201 {object} BucketDTO
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets [post]
func (h *Handler) CreateBucket(c *gin.Context) {
	setResponseHeaders(c)
	var req CreateBucketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if err := validator.ValidateBucketName(req.Bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := h.service.CreateBucket(ctx, req.Bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, &BucketDTO{
		Name:      req.Bucket,
		CreatedAt: time.Now().UTC(),
	})
}

// ListBuckets godoc
// @Summary List storage buckets
// @Description Retrieves all buckets available in the current workspace
// @Tags buckets
// @Produce json
// @Security BearerAuth
// @Success 200 {array} BucketDTO
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets [get]
func (h *Handler) ListBuckets(c *gin.Context) {
	setResponseHeaders(c)
	ctx := c.Request.Context()
	buckets, err := h.service.ListBuckets(ctx)
	if err != nil {
		h.handleStorageError(c, err)
		return
	}

	dtoList := ToBucketDTOList(buckets)
	if dtoList == nil {
		dtoList = []*BucketDTO{}
	}

	response.Success(c, http.StatusOK, dtoList)
}

// DeleteBucket godoc
// @Summary Delete a storage bucket
// @Description Deletes an empty storage bucket
// @Tags buckets
// @Security BearerAuth
// @Param bucket path string true "Bucket name"
// @Success 204 "No Content"
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets/{bucket} [delete]
func (h *Handler) DeleteBucket(c *gin.Context) {
	setResponseHeaders(c)
	bucket := c.Param("bucket")
	if err := validator.ValidateBucketName(bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := h.service.DeleteBucket(ctx, bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UploadObject godoc
// @Summary Stream upload an object
// @Description Uploads a file to the specified bucket using multipart/form-data streaming
// @Tags objects
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Param bucket path string true "Bucket name"
// @Param key query string false "Object key (optional, defaults to uploaded filename)"
// @Param file formData file true "File to upload"
// @Success 201 {object} UploadResponseDTO
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets/{bucket}/objects [post]
func (h *Handler) UploadObject(c *gin.Context) {
	setResponseHeaders(c)
	bucket := c.Param("bucket")
	if err := validator.ValidateBucketName(bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_multipart", "Request is not a valid multipart/form-data stream.")
		return
	}

	var key string
	if k := c.Query("key"); k != "" {
		key = k
	}

	ctx := c.Request.Context()

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			response.Error(c, http.StatusBadRequest, "invalid_multipart", "Failed to read multipart stream part.")
			return
		}

		if part.FormName() == "key" && key == "" {
			keyBytes, _ := io.ReadAll(io.LimitReader(part, 1024))
			key = strings.TrimSpace(string(keyBytes))
			continue
		}

		if part.FormName() == "file" {
			if key == "" {
				key = part.FileName()
			}
			if err := validator.ValidateUpload(key, 0, h.maxUploadSize); err != nil {
				h.handleStorageError(c, err)
				return
			}

			// Wrap part in LimitReader to enforce MAX_UPLOAD_SIZE_MB during streaming
			limitReader := io.LimitReader(part, h.maxUploadSize+1)

			meta := &ObjectMetadata{
				ContentType: part.Header.Get("Content-Type"),
			}
			if meta.ContentType == "" {
				meta.ContentType = "application/octet-stream"
			}

			res, err := h.service.UploadObject(ctx, bucket, key, limitReader, -1, meta)
			if err != nil {
				h.handleStorageError(c, err)
				return
			}

			response.Success(c, http.StatusCreated, ToUploadResponseDTO(res))
			return
		}
	}

	response.Error(c, http.StatusBadRequest, "missing_file", "No file field found in multipart upload request.")
}

// ListObjects godoc
// @Summary List objects in a bucket
// @Description Retrieves a list of objects in a bucket, optionally filtered by prefix
// @Tags objects
// @Produce json
// @Security BearerAuth
// @Param bucket path string true "Bucket name"
// @Param prefix query string false "Key prefix filter"
// @Success 200 {array} ObjectDTO
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets/{bucket}/objects [get]
func (h *Handler) ListObjects(c *gin.Context) {
	setResponseHeaders(c)
	bucket := c.Param("bucket")
	if err := validator.ValidateBucketName(bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}

	prefix := c.Query("prefix")
	ctx := c.Request.Context()
	objects, err := h.service.ListObjects(ctx, bucket, prefix)
	if err != nil {
		h.handleStorageError(c, err)
		return
	}

	dtoList := ToObjectDTOList(objects)
	if dtoList == nil {
		dtoList = []*ObjectDTO{}
	}

	response.Success(c, http.StatusOK, dtoList)
}

// DownloadObject godoc
// @Summary Stream download an object
// @Description Streams an object's contents from the storage provider
// @Tags objects
// @Produce application/octet-stream
// @Security BearerAuth
// @Param bucket path string true "Bucket name"
// @Param key path string true "Object key"
// @Success 200 {file} file
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets/{bucket}/objects/{key} [get]
func (h *Handler) DownloadObject(c *gin.Context) {
	setResponseHeaders(c)
	bucket := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")
	if err := validator.ValidateBucketName(bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}
	if err := validator.ValidateObjectKey(key); err != nil {
		h.handleStorageError(c, err)
		return
	}

	ctx := c.Request.Context()

	// Check metadata / headers first
	if objs, err := h.service.ListObjects(ctx, bucket, key); err == nil && len(objs) > 0 {
		for _, o := range objs {
			if o.Key == key {
				if o.ETag != "" {
					c.Header("ETag", o.ETag)
				}
				if !o.LastModified.IsZero() {
					c.Header("Last-Modified", o.LastModified.UTC().Format(http.TimeFormat))
				}
				c.Header("Cache-Control", "public, max-age=3600")
				if o.Size > 0 {
					c.Header("Content-Length", fmt.Sprintf("%d", o.Size))
				}
				if o.Metadata != nil && o.Metadata.ContentType != "" {
					c.Header("Content-Type", o.Metadata.ContentType)
				} else {
					c.Header("Content-Type", "application/octet-stream")
				}
				break
			}
		}
	} else {
		c.Header("Content-Type", "application/octet-stream")
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(key)))

	rc, err := h.service.DownloadObject(ctx, bucket, key)
	if err != nil {
		h.handleStorageError(c, err)
		return
	}
	defer rc.Close()

	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, rc)
}

// DeleteObject godoc
// @Summary Delete an object
// @Description Deletes an object from the storage bucket
// @Tags objects
// @Security BearerAuth
// @Param bucket path string true "Bucket name"
// @Param key path string true "Object key"
// @Success 204 "No Content"
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets/{bucket}/objects/{key} [delete]
func (h *Handler) DeleteObject(c *gin.Context) {
	setResponseHeaders(c)
	bucket := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")
	if err := validator.ValidateBucketName(bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}
	if err := validator.ValidateObjectKey(key); err != nil {
		h.handleStorageError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := h.service.DeleteObject(ctx, bucket, key); err != nil {
		h.handleStorageError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ObjectExists godoc
// @Summary Check if an object exists
// @Description Returns HTTP 200 if the object exists in the bucket, 404 otherwise
// @Tags objects
// @Security BearerAuth
// @Param bucket path string true "Bucket name"
// @Param key path string true "Object key"
// @Success 200 "OK"
// @Failure 401 {object} response.APIResponse
// @Failure 404 "Not Found"
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/storage/buckets/{bucket}/objects/{key} [head]
func (h *Handler) ObjectExists(c *gin.Context) {
	setResponseHeaders(c)
	bucket := c.Param("bucket")
	key := strings.TrimPrefix(c.Param("key"), "/")
	if err := validator.ValidateBucketName(bucket); err != nil {
		h.handleStorageError(c, err)
		return
	}
	if err := validator.ValidateObjectKey(key); err != nil {
		h.handleStorageError(c, err)
		return
	}

	ctx := c.Request.Context()
	exists, err := h.service.ObjectExists(ctx, bucket, key)
	if err != nil {
		h.handleStorageError(c, err)
		return
	}
	if !exists {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
