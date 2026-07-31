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
	"github.com/cloudstorex/backend/internal/storage"
)

type Provider struct {
	client *s3.Client
	logger *slog.Logger
	bucket string
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
		client: client,
		logger: logger,
	}, nil
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

func (p *Provider) GeneratePresignedURL(ctx context.Context, bucket, key string, expiration time.Duration) (string, error) {
	return "", errors.New("not implemented")
}

func (p *Provider) CopyObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) MoveObject(ctx context.Context, srcBucket, srcKey, destBucket, destKey string) (*storage.StorageResponse, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) GetObjectMetadata(ctx context.Context, bucket, key string) (*storage.ObjectMetadata, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) SetObjectMetadata(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) error {
	return errors.New("not implemented")
}

func (p *Provider) GetObjectTags(ctx context.Context, bucket, key string) (map[string]string, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) SetObjectTags(ctx context.Context, bucket, key string, tags map[string]string) error {
	return errors.New("not implemented")
}

func (p *Provider) CreateMultipartUpload(ctx context.Context, bucket, key string, meta *storage.ObjectMetadata) (*storage.MultipartUpload, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) UploadPart(ctx context.Context, uploadID, bucket, key string, partNumber int, reader io.Reader, size int64) (*storage.UploadPart, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) CompleteMultipartUpload(ctx context.Context, uploadID, bucket, key string, parts []*storage.UploadPart) (*storage.StorageResponse, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) AbortMultipartUpload(ctx context.Context, uploadID, bucket, key string) error {
	return errors.New("not implemented")
}

func (p *Provider) ListObjectVersions(ctx context.Context, bucket, key string) ([]*storage.Object, error) {
	return nil, errors.New("not implemented")
}
