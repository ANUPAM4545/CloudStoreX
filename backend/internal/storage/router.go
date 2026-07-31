package storage

import (
	"context"
	"io"
	"time"

	"github.com/cloudstorex/backend/internal/policy/engine"
	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/google/uuid"
)

// ProviderRegistry defines the lookup interface required by the Storage Router
// to retrieve a StorageProvider by ID without creating a cyclic dependency on the provider package.
type ProviderRegistry interface {
	Get(id string) (StorageProvider, error)
}

// Router defines the central storage routing contract.
// All storage operations must pass through this router.
type Router interface {
	StorageProvider
}

type defaultRouter struct {
	registry  ProviderRegistry
	engine    engine.PolicyEngine
}

// NewRouter creates a new central Storage Router.
func NewRouter(registry ProviderRegistry, policyEngine engine.PolicyEngine) Router {
	return &defaultRouter{
		registry:  registry,
		engine:    policyEngine,
	}
}

// resolveProvider determines the target provider ID using the policy engine
// and retrieves the provider instance from the registry.
func (r *defaultRouter) resolveProvider(ctx context.Context, bucket, key string, size int64, op string) (StorageProvider, string, error) {
	evalCtx := &model.EvaluationContext{
		Bucket:    bucket,
		ObjectKey: key,
		Size:      size,
	}

	// Propagate ID context if available
	if val := GetContextValue(ctx, CtxKeyWorkspaceID); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			evalCtx.WorkspaceID = id
		}
	}

	providerID, _, err := r.engine.Resolve(ctx, evalCtx, op)
	if err != nil {
		return nil, "", NewDomainError("resolve_provider", "", bucket, key, err)
	}

	provider, err := r.registry.Get(providerID)
	if err != nil {
		return nil, providerID, NewDomainError("get_provider", providerID, bucket, key, err)
	}

	return provider, providerID, nil
}

func (r *defaultRouter) Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *ObjectMetadata) (*StorageResponse, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, size, "Upload")
	if err != nil {
		return nil, err
	}
	res, err := prov.Upload(ctx, bucket, key, reader, size, meta)
	if err != nil {
		return nil, NewDomainError("upload", providerID, bucket, key, err)
	}
	return res, nil
}

func (r *defaultRouter) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "Download")
	if err != nil {
		return nil, err
	}
	res, err := prov.Download(ctx, bucket, key)
	if err != nil {
		return nil, NewDomainError("download", providerID, bucket, key, err)
	}
	return res, nil
}

func (r *defaultRouter) Delete(ctx context.Context, bucket, key string) error {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "Delete")
	if err != nil {
		return err
	}
	if err := prov.Delete(ctx, bucket, key); err != nil {
		return NewDomainError("delete", providerID, bucket, key, err)
	}
	return nil
}

func (r *defaultRouter) Exists(ctx context.Context, bucket, key string) (bool, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "Exists")
	if err != nil {
		return false, err
	}
	exists, err := prov.Exists(ctx, bucket, key)
	if err != nil {
		return false, NewDomainError("exists", providerID, bucket, key, err)
	}
	return exists, nil
}

func (r *defaultRouter) ListObjects(ctx context.Context, bucket, prefix string) ([]*Object, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, prefix, 0, "ListObjects")
	if err != nil {
		return nil, err
	}
	objs, err := prov.ListObjects(ctx, bucket, prefix)
	if err != nil {
		return nil, NewDomainError("list_objects", providerID, bucket, prefix, err)
	}
	return objs, nil
}

func (r *defaultRouter) CreateBucket(ctx context.Context, bucket string) error {
	prov, providerID, err := r.resolveProvider(ctx, bucket, "", 0, "CreateBucket")
	if err != nil {
		return err
	}
	if err := prov.CreateBucket(ctx, bucket); err != nil {
		return NewDomainError("create_bucket", providerID, bucket, "", err)
	}
	return nil
}

func (r *defaultRouter) DeleteBucket(ctx context.Context, bucket string) error {
	prov, providerID, err := r.resolveProvider(ctx, bucket, "", 0, "DeleteBucket")
	if err != nil {
		return err
	}
	if err := prov.DeleteBucket(ctx, bucket); err != nil {
		return NewDomainError("delete_bucket", providerID, bucket, "", err)
	}
	return nil
}

func (r *defaultRouter) ListBuckets(ctx context.Context) ([]*Bucket, error) {
	prov, providerID, err := r.resolveProvider(ctx, "", "", 0, "ListBuckets")
	if err != nil {
		return nil, err
	}
	buckets, err := prov.ListBuckets(ctx)
	if err != nil {
		return nil, NewDomainError("list_buckets", providerID, "", "", err)
	}
	return buckets, nil
}

func (r *defaultRouter) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "GeneratePresignedURL")
	if err != nil {
		return "", err
	}
	url, err := prov.GeneratePresignedURL(ctx, bucket, key, expiration)
	if err != nil {
		return "", NewDomainError("presigned_url", providerID, bucket, key, err)
	}
	return url, nil
}

func (r *defaultRouter) CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*StorageResponse, error) {
	prov, providerID, err := r.resolveProvider(ctx, srcBucket, srcKey, 0, "CopyObject")
	if err != nil {
		return nil, err
	}
	res, err := prov.CopyObject(ctx, srcBucket, srcKey, destBucket, destKey)
	if err != nil {
		return nil, NewDomainError("copy_object", providerID, srcBucket, srcKey, err)
	}
	return res, nil
}

func (r *defaultRouter) MoveObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*StorageResponse, error) {
	prov, providerID, err := r.resolveProvider(ctx, srcBucket, srcKey, 0, "MoveObject")
	if err != nil {
		return nil, err
	}
	res, err := prov.MoveObject(ctx, srcBucket, srcKey, destBucket, destKey)
	if err != nil {
		return nil, NewDomainError("move_object", providerID, srcBucket, srcKey, err)
	}
	return res, nil
}

func (r *defaultRouter) GetObjectMetadata(ctx context.Context, bucket, key string) (*ObjectMetadata, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "GetObjectMetadata")
	if err != nil {
		return nil, err
	}
	meta, err := prov.GetObjectMetadata(ctx, bucket, key)
	if err != nil {
		return nil, NewDomainError("get_metadata", providerID, bucket, key, err)
	}
	return meta, nil
}

func (r *defaultRouter) SetObjectMetadata(ctx context.Context, bucket, key string, meta *ObjectMetadata) error {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "SetObjectMetadata")
	if err != nil {
		return err
	}
	if err := prov.SetObjectMetadata(ctx, bucket, key, meta); err != nil {
		return NewDomainError("set_metadata", providerID, bucket, key, err)
	}
	return nil
}

func (r *defaultRouter) GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "GetObjectTags")
	if err != nil {
		return nil, err
	}
	tags, err := prov.GetObjectTags(ctx, bucket, key)
	if err != nil {
		return nil, NewDomainError("get_tags", providerID, bucket, key, err)
	}
	return tags, nil
}

func (r *defaultRouter) SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "SetObjectTags")
	if err != nil {
		return err
	}
	if err := prov.SetObjectTags(ctx, bucket, key, tags); err != nil {
		return NewDomainError("set_tags", providerID, bucket, key, err)
	}
	return nil
}

func (r *defaultRouter) CreateMultipartUpload(ctx context.Context, bucket, key string, meta *ObjectMetadata) (*MultipartUpload, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "CreateMultipartUpload")
	if err != nil {
		return nil, err
	}
	res, err := prov.CreateMultipartUpload(ctx, bucket, key, meta)
	if err != nil {
		return nil, NewDomainError("create_multipart", providerID, bucket, key, err)
	}
	return res, nil
}

func (r *defaultRouter) UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*UploadPart, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, size, "UploadPart")
	if err != nil {
		return nil, err
	}
	res, err := prov.UploadPart(ctx, uploadID, bucket, key, partNumber, reader, size)
	if err != nil {
		return nil, NewDomainError("upload_part", providerID, bucket, key, err)
	}
	return res, nil
}

func (r *defaultRouter) CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*UploadPart) (*StorageResponse, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "CompleteMultipartUpload")
	if err != nil {
		return nil, err
	}
	res, err := prov.CompleteMultipartUpload(ctx, uploadID, bucket, key, parts)
	if err != nil {
		return nil, NewDomainError("complete_multipart", providerID, bucket, key, err)
	}
	return res, nil
}

func (r *defaultRouter) AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "AbortMultipartUpload")
	if err != nil {
		return err
	}
	if err := prov.AbortMultipartUpload(ctx, uploadID, bucket, key); err != nil {
		return NewDomainError("abort_multipart", providerID, bucket, key, err)
	}
	return nil
}

func (r *defaultRouter) ListObjectVersions(ctx context.Context, bucket, key string) ([]*Object, error) {
	prov, providerID, err := r.resolveProvider(ctx, bucket, key, 0, "ListObjectVersions")
	if err != nil {
		return nil, err
	}
	versions, err := prov.ListObjectVersions(ctx, bucket, key)
	if err != nil {
		return nil, NewDomainError("list_versions", providerID, bucket, key, err)
	}
	return versions, nil
}
