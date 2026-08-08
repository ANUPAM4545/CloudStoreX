package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/service"
	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/cloudstorex/backend/internal/quota"
	"github.com/cloudstorex/backend/internal/shared/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

// Service defines the application-level storage contract.
type Service interface {
	UploadObject(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error)
	DownloadObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucket, key string) error
	ObjectExists(ctx context.Context, bucket, key string) (bool, error)
	ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error)
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error
	ListBuckets(ctx context.Context) ([]*Bucket, error)
	GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error)
	CopyObject(ctx context.Context, bucket, srcKey, destKey string) error
	MoveObject(ctx context.Context, bucket, srcKey, destKey string) error

	SearchObjects(ctx context.Context, query dto.SearchQuery) ([]dto.ObjectMetaDTO, int64, error)
	GetObjectByID(ctx context.Context, id string) (*dto.ObjectMetaDTO, error)
	GetObjectMetadata(ctx context.Context, id string) (map[string]string, error)
	ListObjectVersions(ctx context.Context, id string) ([]dto.ObjectVersionDTO, error)
	RestoreObject(ctx context.Context, id string) error
	TagObject(ctx context.Context, id string, tags map[string]string) error
	UntagObject(ctx context.Context, id string, keys []string) error
}

type defaultService struct {
	router   Router
	metadata service.MetadataService
	quota    quota.Service
	log      *slog.Logger
}

// NewService creates a new high-level Storage Service.
func NewService(router Router, metadata service.MetadataService, quota quota.Service) Service {
	l := logger.Log
	if l == nil {
		l = slog.Default()
	}
	return &defaultService{
		router:   router,
		metadata: metadata,
		quota:    quota,
		log:      l,
	}
}

func (s *defaultService) logOperation(ctx context.Context, op, bucket, key string, start time.Time, err error) {
	if s.log == nil {
		s.log = slog.Default()
	}
	duration := time.Since(start)
	reqID := GetContextValue(ctx, CtxKeyRequestID)
	workspaceID := GetContextValue(ctx, CtxKeyWorkspaceID)
	traceID := tracing.ExtractTraceID(ctx)
	spanID := tracing.ExtractSpanID(ctx)

	level := slog.LevelInfo
	result := "success"
	var errStr string

	if err != nil {
		level = slog.LevelError
		result = "error"
		errStr = err.Error()
	}

	s.log.Log(ctx, level, "Storage Operation",
		slog.String("op", op),
		slog.String("bucket", bucket),
		slog.String("key", key),
		slog.String("request_id", reqID),
		slog.String("workspace_id", workspaceID),
		slog.String("trace_id", traceID),
		slog.String("span_id", spanID),
		slog.Duration("duration", duration),
		slog.String("result", result),
		slog.String("error", errStr),
	)
}

func (s *defaultService) UploadObject(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error) {
	ctx, span := tracing.StartChildSpan(ctx, "StorageService.UploadObject")
	defer span.End()
	span.SetAttributes(
		attribute.String("storage.bucket", bucket),
		attribute.String("storage.operation", "upload"),
		attribute.Int64("storage.size_bytes", size),
	)

	start := time.Now()
	
	workspaceIDStr := GetContextValue(ctx, CtxKeyWorkspaceID)

	// 1. Quota Check
	if s.quota != nil {
		if err := s.quota.CheckQuota(ctx, workspaceIDStr, size); err != nil {
			s.logOperation(ctx, "UploadObject (Quota Check)", bucket, key, start, err)
			return nil, err
		}
	}
	
	// Set Evaluation Context
	mimeType := "application/octet-stream"
	var tags map[string]string
	var custom map[string]string
	if meta != nil {
		if meta.ContentType != "" {
			mimeType = meta.ContentType
		}
		tags = meta.Tags
		custom = meta.Custom
	}
	
	evalCtx := context.WithValue(ctx, CtxKeyEvaluationMimeType, mimeType)
	evalCtx = context.WithValue(evalCtx, CtxKeyEvaluationTags, tags)
	
	// Upload to Provider first
	res, err := s.router.Upload(evalCtx, bucket, key, reader, size, meta)
	if err != nil {
		s.logOperation(ctx, "UploadObject (Provider)", bucket, key, start, err)
		return nil, err
	}

	// Fetch Bucket from Metadata
	bucketMeta, err := s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err != nil {
		// Rollback physical upload since metadata lookup failed
		_ = s.router.Delete(context.Background(), bucket, key)
		s.logOperation(ctx, "UploadObject (Metadata Rollback)", bucket, key, start, err)
		return nil, fmt.Errorf("metadata bucket lookup failed: %w", err)
	}
	
	var ownerID *uuid.UUID
	userIDStr := GetContextValue(ctx, CtxKeyUserID)
	if userIDStr != "" {
		uid, _ := uuid.Parse(userIDStr)
		ownerID = &uid
	}
	
	// Ensure ETag is grabbed from Provider response
	etag := ""
	if res != nil {
		etag = res.ETag
	}

	// Create Metadata in Postgres
	_, err = s.metadata.CreateObjectMetadata(ctx, workspaceIDStr, bucketMeta.ID, key, key, size, mimeType, etag, bucketMeta.ProviderID, ownerID, tags, custom)
	if err != nil {
		// Rollback physical upload since metadata persistence failed
		_ = s.router.Delete(context.Background(), bucket, key)
		s.logOperation(ctx, "UploadObject (Metadata Persistence Rollback)", bucket, key, start, err)
		return nil, fmt.Errorf("failed to persist metadata: %w", err)
	}

	s.logOperation(ctx, "UploadObject", bucket, key, start, nil)
	return res, nil
}

func (s *defaultService) DownloadObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	ctx, span := tracing.StartChildSpan(ctx, "StorageService.DownloadObject")
	defer span.End()
	span.SetAttributes(
		attribute.String("storage.bucket", bucket),
		attribute.String("storage.operation", "download"),
	)

	start := time.Now()
	
	workspaceIDStr := GetContextValue(ctx, CtxKeyWorkspaceID)
	bucketMeta, err := s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err != nil {
		s.logOperation(ctx, "DownloadObject (Metadata Check)", bucket, key, start, err)
		return nil, err
	}
	
	// Ensure object exists in metadata first
	_, err = s.metadata.FindObjectByKey(ctx, bucketMeta.ID.String(), key)
	if err != nil {
		s.logOperation(ctx, "DownloadObject (Metadata Not Found)", bucket, key, start, err)
		return nil, err
	}

	res, err := s.router.Download(ctx, bucket, key)
	s.logOperation(ctx, "DownloadObject", bucket, key, start, err)
	return res, err
}

func (s *defaultService) DeleteObject(ctx context.Context, bucket, key string) error {
	ctx, span := tracing.StartChildSpan(ctx, "StorageService.DeleteObject")
	defer span.End()
	span.SetAttributes(
		attribute.String("storage.bucket", bucket),
		attribute.String("storage.operation", "delete"),
	)

	start := time.Now()
	
	workspaceIDStr := GetContextValue(ctx, CtxKeyWorkspaceID)
	bucketMeta, err := s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err != nil {
		return err
	}
	
	objMeta, err := s.metadata.FindObjectByKey(ctx, bucketMeta.ID.String(), key)
	if err != nil {
		return err
	}

	// Compliance Checks
	if objMeta.LegalHold {
		s.logOperation(ctx, "DeleteObject (Blocked by Legal Hold)", bucket, key, start, fmt.Errorf("object is under legal hold"))
		return fmt.Errorf("cannot delete object under legal hold")
	}

	if objMeta.RetainUntil != nil && objMeta.RetainUntil.After(time.Now()) {
		s.logOperation(ctx, "DeleteObject (Blocked by Retention)", bucket, key, start, fmt.Errorf("object is retained until %v", objMeta.RetainUntil))
		return fmt.Errorf("cannot delete object, retained until %v", objMeta.RetainUntil)
	}

	// Delete from Provider
	err = s.router.Delete(ctx, bucket, key)
	if err != nil {
		s.logOperation(ctx, "DeleteObject (Provider)", bucket, key, start, err)
		return err
	}
	
	// Soft Delete Metadata
	err = s.metadata.SoftDeleteObject(ctx, workspaceIDStr, objMeta.ID.String())
	s.logOperation(ctx, "DeleteObject", bucket, key, start, err)
	return err
}

func (s *defaultService) ObjectExists(ctx context.Context, bucket, key string) (bool, error) {
	start := time.Now()
	
	workspaceIDStr := GetContextValue(ctx, CtxKeyWorkspaceID)
	bucketMeta, err := s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err != nil {
		return false, err
	}
	
	_, err = s.metadata.FindObjectByKey(ctx, bucketMeta.ID.String(), key)
	if err != nil {
		s.logOperation(ctx, "ObjectExists", bucket, key, start, nil)
		return false, nil // Not found in metadata
	}

	s.logOperation(ctx, "ObjectExists", bucket, key, start, nil)
	return true, nil
}

func (s *defaultService) ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error) {
	start := time.Now()
	
	workspaceIDStr := GetContextValue(ctx, CtxKeyWorkspaceID)
	bucketMeta, err := s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err != nil {
		return nil, err
	}
	
	// Fetch from Metadata Catalog
	results, _, err := s.metadata.SearchObjects(ctx, dto.SearchQuery{
		BucketID: bucketMeta.ID.String(),
		Prefix:   prefix,
		Limit:    1000,
	})
	
	if err != nil {
		s.logOperation(ctx, "ListObjects (Metadata)", bucket, prefix, start, err)
		return nil, err
	}

	var objs []*Object
	for _, res := range results {
		meta := &ObjectMetadata{
			ContentType: res.MimeType,
			Tags:        res.Tags,
			Custom:      res.Metadata,
			VersionID:   res.VersionID,
		}
		
		objs = append(objs, &Object{
			Key:          res.ObjectKey,
			Bucket:       bucket,
			Size:         res.SizeBytes,
			ETag:         res.ETag,
			LastModified: res.CreatedAt,
			Metadata:     meta,
		})
	}
	
	s.logOperation(ctx, "ListObjects", bucket, prefix, start, nil)
	return objs, nil
}

func (s *defaultService) CreateBucket(ctx context.Context, bucket string) error {
	start := time.Now()
	err := s.router.CreateBucket(ctx, bucket)
	s.logOperation(ctx, "CreateBucket", bucket, "", start, err)
	return err
}

func (s *defaultService) DeleteBucket(ctx context.Context, bucket string) error {
	start := time.Now()
	err := s.router.DeleteBucket(ctx, bucket)
	s.logOperation(ctx, "DeleteBucket", bucket, "", start, err)
	return err
}

func (s *defaultService) ListBuckets(ctx context.Context) ([]*Bucket, error) {
	start := time.Now()
	buckets, err := s.router.ListBuckets(ctx)
	s.logOperation(ctx, "ListBuckets", "", "", start, err)
	return buckets, err
}

func (s *defaultService) GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	start := time.Now()
	url, err := s.router.GeneratePresignedUploadURL(ctx, bucket, key, expiration)
	s.logOperation(ctx, "GeneratePresignedUploadURL", bucket, key, start, err)
	return url, err
}

func (s *defaultService) GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	start := time.Now()
	url, err := s.router.GeneratePresignedDownloadURL(ctx, bucket, key, expiration)
	s.logOperation(ctx, "GeneratePresignedDownloadURL", bucket, key, start, err)
	return url, err
}

func (s *defaultService) SearchObjects(ctx context.Context, query dto.SearchQuery) ([]dto.ObjectMetaDTO, int64, error) {
	start := time.Now()
	res, total, err := s.metadata.SearchObjects(ctx, query)
	s.logOperation(ctx, "SearchObjects", "", "", start, err)
	return res, total, err
}

func (s *defaultService) GetObjectByID(ctx context.Context, id string) (*dto.ObjectMetaDTO, error) {
	start := time.Now()
	obj, err := s.metadata.FindObjectByID(ctx, id)
	if err != nil {
		s.logOperation(ctx, "GetObjectByID", "", id, start, err)
		return nil, err
	}
	
	t := make(map[string]string)
	for _, tag := range obj.Tags {
		t[tag.Key] = tag.Value
	}
	m := make(map[string]string)
	for _, meta := range obj.Metadata {
		m[meta.Key] = meta.Value
	}
	
	dtoObj := &dto.ObjectMetaDTO{
		ID:                obj.ID.String(),
		BucketID:          obj.BucketID.String(),
		ObjectKey:         obj.ObjectKey,
		ProviderObjectKey: obj.ProviderObjectKey,
		SizeBytes:         obj.SizeBytes,
		MimeType:          obj.MimeType,
		ETag:              obj.ETag,
		ProviderID:        obj.ProviderID,
		StorageClass:      obj.StorageClass,
		Status:            obj.Status,
		Tags:              t,
		Metadata:          m,
		CreatedAt:         obj.CreatedAt,
		UpdatedAt:         obj.UpdatedAt,
	}

	s.logOperation(ctx, "GetObjectByID", "", id, start, nil)
	return dtoObj, nil
}

func (s *defaultService) GetObjectMetadata(ctx context.Context, id string) (map[string]string, error) {
	start := time.Now()
	res, err := s.metadata.GetObjectMetadata(ctx, id)
	s.logOperation(ctx, "GetObjectMetadata", "", id, start, err)
	return res, err
}

func (s *defaultService) ListObjectVersions(ctx context.Context, id string) ([]dto.ObjectVersionDTO, error) {
	start := time.Now()
	res, err := s.metadata.ListObjectVersions(ctx, id)
	s.logOperation(ctx, "ListObjectVersions", "", id, start, err)
	return res, err
}

func (s *defaultService) RestoreObject(ctx context.Context, id string) error {
	start := time.Now()
	err := s.metadata.RestoreObject(ctx, id)
	s.logOperation(ctx, "RestoreObject", "", id, start, err)
	return err
}

func (s *defaultService) TagObject(ctx context.Context, id string, tags map[string]string) error {
	start := time.Now()
	err := s.metadata.TagObject(ctx, id, tags)
	s.logOperation(ctx, "TagObject", "", id, start, err)
	return err
}

func (s *defaultService) UntagObject(ctx context.Context, id string, keys []string) error {
	start := time.Now()
	err := s.metadata.UntagObject(ctx, id, keys)
	s.logOperation(ctx, "UntagObject", "", id, start, err)
	return err
}

func (s *defaultService) CopyObject(ctx context.Context, bucket, srcKey, destKey string) error {
	ctx, span := tracing.StartChildSpan(ctx, "StorageService.CopyObject")
	defer span.End()
	span.SetAttributes(
		attribute.String("storage.bucket", bucket),
		attribute.String("storage.src_key", srcKey),
		attribute.String("storage.dest_key", destKey),
	)

	start := time.Now()
	
	workspaceIDStr := GetContextValue(ctx, CtxKeyWorkspaceID)
	evalCtx := ctx
	
	// Pre-fetch metadata to populate EvaluationContext for CopyObject
	bucketMeta, err := s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err == nil {
		if srcMeta, err := s.metadata.FindObjectByKey(ctx, bucketMeta.ID.String(), srcKey); err == nil {
			evalCtx = context.WithValue(evalCtx, CtxKeyEvaluationMimeType, srcMeta.MimeType)
			
			t := make(map[string]string)
			for _, tag := range srcMeta.Tags {
				t[tag.Key] = tag.Value
			}
			evalCtx = context.WithValue(evalCtx, CtxKeyEvaluationTags, t)
		}
	}

	// Delegate copy to router/provider
	_, err = s.router.CopyObject(evalCtx, bucket, srcKey, bucket, destKey)
	if err != nil {
		s.logOperation(ctx, "CopyObject (Provider)", bucket, srcKey, start, err)
		return err
	}

	// Update metadata: find source, copy metadata to dest (this logic can be further enhanced)
	workspaceIDStr = GetContextValue(ctx, CtxKeyWorkspaceID)
	bucketMeta, err = s.metadata.FindBucketByName(ctx, workspaceIDStr, bucket)
	if err == nil {
		srcMeta, err := s.metadata.FindObjectByKey(ctx, bucketMeta.ID.String(), srcKey)
		if err == nil {
			var ownerID *uuid.UUID
			if srcMeta.OwnerID != nil {
				ownerID = srcMeta.OwnerID
			}
			t := make(map[string]string)
			for _, tag := range srcMeta.Tags {
				t[tag.Key] = tag.Value
			}
			m := make(map[string]string)
			for _, meta := range srcMeta.Metadata {
				m[meta.Key] = meta.Value
			}
			
			_, _ = s.metadata.CreateObjectMetadata(ctx, workspaceIDStr, bucketMeta.ID, destKey, destKey, srcMeta.SizeBytes, srcMeta.MimeType, srcMeta.ETag, srcMeta.ProviderID, ownerID, t, m)
		}
	}

	s.logOperation(ctx, "CopyObject", bucket, srcKey, start, nil)
	return nil
}

func (s *defaultService) MoveObject(ctx context.Context, bucket, srcKey, destKey string) error {
	ctx, span := tracing.StartChildSpan(ctx, "StorageService.MoveObject")
	defer span.End()
	span.SetAttributes(
		attribute.String("storage.bucket", bucket),
		attribute.String("storage.src_key", srcKey),
		attribute.String("storage.dest_key", destKey),
	)

	start := time.Now()

	// 1. CopyObject
	err := s.CopyObject(ctx, bucket, srcKey, destKey)
	if err != nil {
		s.logOperation(ctx, "MoveObject (Copy)", bucket, srcKey, start, err)
		return err
	}

	// 2. Delete original object
	err = s.DeleteObject(ctx, bucket, srcKey)
	if err != nil {
		s.logOperation(ctx, "MoveObject (Delete Src)", bucket, srcKey, start, err)
		return err
	}

	s.logOperation(ctx, "MoveObject", bucket, srcKey, start, nil)
	return nil
}

