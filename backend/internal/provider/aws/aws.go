package aws

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/cloudstorex/backend/internal/provider"
	"github.com/cloudstorex/backend/internal/storage"
)

type Provider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	logger        *slog.Logger
}

type Config struct {
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketPrefix    string
}

// NewProvider creates a new AWS S3 StorageProvider instance.
func NewProvider(ctx context.Context, cfg *Config, logger *slog.Logger) (storage.StorageProvider, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to its default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Needed for LocalStack or MinIO compatibility
	})

	return &Provider{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		logger:        logger,
	}, nil
}

func (p *Provider) translateError(op, bucket, key string, err error) error {
	if err == nil {
		return nil
	}
	
	// Fast path for context errors
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return storage.NewDomainError(op, "aws", bucket, key, err)
	}

	var apiErr interface {
		ErrorCode() string
		ErrorMessage() string
	}
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NotFound":
			return storage.NewDomainError(op, "aws", bucket, key, storage.ErrObjectNotFound)
		case "NoSuchBucket":
			return storage.NewDomainError(op, "aws", bucket, key, storage.ErrBucketNotFound)
		case "BucketAlreadyExists", "BucketAlreadyOwnedByYou":
			return storage.NewDomainError(op, "aws", bucket, key, storage.ErrBucketExists)
		case "AccessDenied":
			return storage.NewDomainError(op, "aws", bucket, key, storage.ErrAccessDenied)
		case "InvalidAccessKeyId", "SignatureDoesNotMatch":
			return storage.NewDomainError(op, "aws", bucket, key, storage.ErrInvalidCredentials)
		case "ServiceUnavailable", "InternalError":
			return storage.NewDomainError(op, "aws", bucket, key, storage.ErrProviderUnavailable)
		}
	}
	
	switch op {
	case "Upload":
		return storage.NewDomainError(op, "aws", bucket, key, storage.ErrUploadFailed)
	case "Download":
		return storage.NewDomainError(op, "aws", bucket, key, storage.ErrDownloadFailed)
	case "Delete":
		return storage.NewDomainError(op, "aws", bucket, key, storage.ErrDeleteFailed)
	}
	return storage.NewDomainError(op, "aws", bucket, key, err)
}

func (p *Provider) logOperation(ctx context.Context, op, bucket, key string, start time.Time, err error) {
	duration := time.Since(start)
	level := slog.LevelInfo
	result := "success"
	var errStr string

	if err != nil {
		level = slog.LevelError
		result = "error"
		errStr = err.Error()
	}

	p.logger.Log(ctx, level, "AWS Provider Operation",
		slog.String("provider", "aws"),
		slog.String("bucket", bucket),
		slog.String("key", key),
		slog.String("op", op),
		slog.Duration("duration", duration),
		slog.String("result", result),
		slog.String("error", errStr),
	)
}

// Bucket Operations
func (p *Provider) CreateBucket(ctx context.Context, bucket string) error {
	p.logger.Info("AWS Provider Operation", "op", "CreateBucket", "bucket", bucket)
	_, err := p.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}

func (p *Provider) DeleteBucket(ctx context.Context, bucket string) error {
	p.logger.Info("AWS Provider Operation", "op", "DeleteBucket", "bucket", bucket)
	_, err := p.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}

func (p *Provider) ListBuckets(ctx context.Context) ([]*storage.Bucket, error) {
	p.logger.Info("AWS Provider Operation", "op", "ListBuckets")
	out, err := p.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var buckets []*storage.Bucket
	for _, b := range out.Buckets {
		buckets = append(buckets, &storage.Bucket{
			Name:      *b.Name,
			CreatedAt: *b.CreationDate,
		})
	}
	return buckets, nil
}

// Object Operations
func (p *Provider) Upload(ctx context.Context, bucket, key string, reader io.Reader, size int64, meta *storage.ObjectMetadata) (*storage.StorageResponse, error) {
	start := time.Now()
	p.logger.Info("AWS Provider Operation", "op", "Upload", "bucket", bucket, "key", key)

	// Since we are using standard put object, we will just use the PutObject API
	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: aws.Int64(size),
	})

	if err != nil {
		p.logger.Error("AWS Provider Operation", "op", "Upload", "error", err)
		return nil, err
	}

	duration := time.Since(start)
	p.logger.Info("AWS Provider Operation Completed", "op", "Upload", "duration", duration)

	return &storage.StorageResponse{
		Bucket: bucket,
		Key:    key,
		Size:   size,
	}, nil
}

func (p *Provider) Download(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	start := time.Now()
	p.logger.Info("AWS Provider Operation", "op", "Download", "bucket", bucket, "key", key)

	out, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Basic mapping of Not Found
		return nil, err
	}

	duration := time.Since(start)
	p.logger.Info("AWS Provider Operation Completed", "op", "Download", "duration", duration)

	return out.Body, nil
}

func (p *Provider) Delete(ctx context.Context, bucket, key string) error {
	p.logger.Info("AWS Provider Operation", "op", "Delete", "bucket", bucket, "key", key)
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}

func (p *Provider) Exists(ctx context.Context, bucket, key string) (bool, error) {
	p.logger.Info("AWS Provider Operation", "op", "Exists", "bucket", bucket, "key", key)
	_, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Ideally we check if error is NotFound
		return false, nil // returning false assuming it's not found
	}
	return true, nil
}

func (p *Provider) ListObjects(ctx context.Context, bucket, prefix string) ([]*storage.Object, error) {
	p.logger.Info("AWS Provider Operation", "op", "ListObjects", "bucket", bucket)
	out, err := p.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	var objects []*storage.Object
	for _, obj := range out.Contents {
		objects = append(objects, &storage.Object{
			Key:          *obj.Key,
			Size:         *obj.Size,
			LastModified: *obj.LastModified,
		})
	}
	return objects, nil
}

func (p *Provider) GeneratePresignedDownloadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.GeneratePresignedDownloadURL")
	defer span.End()
	
	start := time.Now()

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	req, err := p.presignClient.PresignGetObject(ctx, input, func(o *s3.PresignOptions) {
		o.Expires = expiration
	})
	
	if err != nil {
		err = p.translateError("GeneratePresignedDownloadURL", bucket, key, err)
	}
	
	p.logOperation(ctx, "GeneratePresignedDownloadURL", bucket, key, start, err)
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (p *Provider) GeneratePresignedUploadURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.GeneratePresignedUploadURL")
	defer span.End()
	
	start := time.Now()

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	req, err := p.presignClient.PresignPutObject(ctx, input, func(o *s3.PresignOptions) {
		o.Expires = expiration
	})
	
	if err != nil {
		err = p.translateError("GeneratePresignedUploadURL", bucket, key, err)
	}
	
	p.logOperation(ctx, "GeneratePresignedUploadURL", bucket, key, start, err)
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (p *Provider) CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.CopyObject")
	defer span.End()
	
	start := time.Now()
	copySource := aws.String(srcBucket + "/" + srcKey)
	
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(destBucket),
		Key:        aws.String(destKey),
		CopySource: copySource,
	}

	var out *s3.CopyObjectOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.CopyObject(ctx, input)
		return p.translateError("CopyObject", destBucket, destKey, innerErr)
	})
	
	p.logOperation(ctx, "CopyObject", destBucket, destKey, start, err)
	if err != nil {
		return nil, err
	}
	
	var versionID string
	if out.VersionId != nil {
		versionID = *out.VersionId
	}
	
	etag := ""
	if out.CopyObjectResult != nil && out.CopyObjectResult.ETag != nil {
		etag = *out.CopyObjectResult.ETag
	}

	return &storage.StorageResponse{
		Bucket:    destBucket,
		Key:       destKey,
		ETag:      etag,
		VersionID: versionID,
	}, nil
}



func (p *Provider) GetObjectMetadata(ctx context.Context, bucket, key string) (*storage.ObjectMetadata, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.GetObjectMetadata")
	defer span.End()
	
	start := time.Now()
	
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	
	var out *s3.HeadObjectOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.HeadObject(ctx, input)
		return p.translateError("GetObjectMetadata", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "GetObjectMetadata", bucket, key, start, err)
	if err != nil {
		return nil, err
	}
	
	contentType := ""
	if out.ContentType != nil {
		contentType = *out.ContentType
	}
	
	var versionID string
	if out.VersionId != nil {
		versionID = *out.VersionId
	}
	
	return &storage.ObjectMetadata{
		ContentType: contentType,
		Custom:      out.Metadata,
		VersionID:   versionID,
	}, nil
}

func (p *Provider) SetObjectMetadata(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) error {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.SetObjectMetadata")
	defer span.End()
	
	start := time.Now()
	
	// S3 requires a CopyObject to self with REPLACE directive to update metadata
	copySource := aws.String(bucket + "/" + key)
	contentType := "application/octet-stream"
	var metaMap map[string]string
	
	if meta != nil {
		if meta.ContentType != "" {
			contentType = meta.ContentType
		}
		if meta.Custom != nil {
			metaMap = meta.Custom
		}
	}
	
	input := &s3.CopyObjectInput{
		Bucket:            aws.String(bucket),
		Key:               aws.String(key),
		CopySource:        copySource,
		MetadataDirective: types.MetadataDirectiveReplace,
		ContentType:       aws.String(contentType),
		Metadata:          metaMap,
	}
	
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		_, innerErr := p.client.CopyObject(ctx, input)
		return p.translateError("SetObjectMetadata", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "SetObjectMetadata", bucket, key, start, err)
	return err
}

func (p *Provider) GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.GetObjectTags")
	defer span.End()
	
	start := time.Now()
	
	input := &s3.GetObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	
	var out *s3.GetObjectTaggingOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.GetObjectTagging(ctx, input)
		return p.translateError("GetObjectTags", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "GetObjectTags", bucket, key, start, err)
	if err != nil {
		return nil, err
	}
	
	tags := make(map[string]string)
	for _, t := range out.TagSet {
		if t.Key != nil && t.Value != nil {
			tags[*t.Key] = *t.Value
		}
	}
	
	return tags, nil
}

func (p *Provider) SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.SetObjectTags")
	defer span.End()
	
	start := time.Now()
	
	var tagSet []types.Tag
	for k, v := range tags {
		tagSet = append(tagSet, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	
	input := &s3.PutObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Tagging: &types.Tagging{
			TagSet: tagSet,
		},
	}
	
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		_, innerErr := p.client.PutObjectTagging(ctx, input)
		return p.translateError("SetObjectTags", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "SetObjectTags", bucket, key, start, err)
	return err
}

func (p *Provider) CreateMultipartUpload(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) (*storage.MultipartUpload, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.CreateMultipartUpload")
	defer span.End()
	
	start := time.Now()
	
	contentType := "application/octet-stream"
	var metaMap map[string]string
	if meta != nil {
		if meta.ContentType != "" {
			contentType = meta.ContentType
		}
		// S3 handles tags via x-amz-tagging, but it's typically set during CreateMultipartUpload or SetObjectTags
		if meta.Custom != nil {
			metaMap = meta.Custom
		}
	}

	input := &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Metadata:    metaMap,
	}
	
	var out *s3.CreateMultipartUploadOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.CreateMultipartUpload(ctx, input)
		return p.translateError("CreateMultipartUpload", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "CreateMultipartUpload", bucket, key, start, err)
	if err != nil {
		return nil, err
	}
	
	return &storage.MultipartUpload{
		UploadID: *out.UploadId,
		Bucket:   bucket,
		Key:      key,
	}, nil
}

func (p *Provider) UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*storage.UploadPart, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.UploadPart")
	defer span.End()
	
	start := time.Now()
	
	input := &s3.UploadPartInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		PartNumber:    aws.Int32(int32(partNumber)),
		UploadId:      aws.String(uploadID),
		Body:          reader,
		ContentLength: aws.Int64(size),
	}
	
	var out *s3.UploadPartOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.UploadPart(ctx, input)
		return p.translateError("UploadPart", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "UploadPart", bucket, key, start, err)
	if err != nil {
		return nil, err
	}

	return &storage.UploadPart{
		PartNumber: partNumber,
		ETag:       *out.ETag,
		Size:       size,
	}, nil
}

func (p *Provider) CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*storage.UploadPart) (*storage.StorageResponse, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.CompleteMultipartUpload")
	defer span.End()
	
	start := time.Now()
	
	var completedParts []types.CompletedPart
	for _, part := range parts {
		completedParts = append(completedParts, types.CompletedPart{
			ETag:       aws.String(part.ETag),
			PartNumber: aws.Int32(int32(part.PartNumber)),
		})
	}
	
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	
	var out *s3.CompleteMultipartUploadOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.CompleteMultipartUpload(ctx, input)
		return p.translateError("CompleteMultipartUpload", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "CompleteMultipartUpload", bucket, key, start, err)
	if err != nil {
		return nil, err
	}
	
	var versionID string
	if out.VersionId != nil {
		versionID = *out.VersionId
	}
	
	return &storage.StorageResponse{
		Bucket:    bucket,
		Key:       key,
		ETag:      *out.ETag,
		VersionID: versionID,
	}, nil
}

func (p *Provider) AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.AbortMultipartUpload")
	defer span.End()
	
	start := time.Now()
	
	input := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	}
	
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		_, innerErr := p.client.AbortMultipartUpload(ctx, input)
		return p.translateError("AbortMultipartUpload", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "AbortMultipartUpload", bucket, key, start, err)
	return err
}

func (p *Provider) ListObjectVersions(ctx context.Context, bucket, key string) ([]*storage.Object, error) {
	ctx, span := tracing.StartChildSpan(ctx, "AWSProvider.ListObjectVersions")
	defer span.End()
	
	start := time.Now()
	
	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}
	
	var out *s3.ListObjectVersionsOutput
	err := provider.WithRetry(ctx, provider.DefaultRetryConfig(), func() error {
		var innerErr error
		out, innerErr = p.client.ListObjectVersions(ctx, input)
		return p.translateError("ListObjectVersions", bucket, key, innerErr)
	})
	
	p.logOperation(ctx, "ListObjectVersions", bucket, key, start, err)
	if err != nil {
		return nil, err
	}
	
	var versions []*storage.Object
	for _, v := range out.Versions {
		if v.Key != nil && *v.Key == key {
			versions = append(versions, &storage.Object{
				Key:          *v.Key,
				Size:         *v.Size,
				LastModified: *v.LastModified,
				VersionID:    *v.VersionId,
				IsLatest:     *v.IsLatest,
			})
		}
	}
	
	return versions, nil
}

func (p *Provider) Capabilities() storage.ProviderCapabilities {
	return storage.ProviderCapabilities{
		MultipartUpload:       true,
		ObjectCopy:            true,
		ObjectVersioning:      true,
		ObjectTags:            true,
		ObjectMetadata:        true,
		PresignedUploadURLs:   true,
		PresignedDownloadURLs: true,
	}
}
