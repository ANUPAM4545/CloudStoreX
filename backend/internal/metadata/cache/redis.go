package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/model"
	"github.com/redis/go-redis/v9"
)

type MetadataCache interface {
	GetObject(ctx context.Context, bucketID, objectKey string) (*model.Object, error)
	SetObject(ctx context.Context, obj *model.Object, ttl time.Duration) error
	InvalidateObject(ctx context.Context, bucketID, objectKey string) error
	
	GetBucket(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error)
	SetBucket(ctx context.Context, bucket *model.Bucket, ttl time.Duration) error
}

type redisMetadataCache struct {
	client *redis.Client
}

func NewRedisMetadataCache(client *redis.Client) MetadataCache {
	return &redisMetadataCache{client: client}
}

func objectCacheKey(bucketID, objectKey string) string {
	return fmt.Sprintf("metadata:object:%s:%s", bucketID, objectKey)
}

func bucketCacheKey(workspaceID, bucketName string) string {
	return fmt.Sprintf("metadata:bucket:%s:%s", workspaceID, bucketName)
}

func (c *redisMetadataCache) GetObject(ctx context.Context, bucketID, objectKey string) (*model.Object, error) {
	key := objectCacheKey(bucketID, objectKey)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var obj model.Object
	if err := json.Unmarshal([]byte(val), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (c *redisMetadataCache) SetObject(ctx context.Context, obj *model.Object, ttl time.Duration) error {
	key := objectCacheKey(obj.BucketID.String(), obj.ObjectKey)
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *redisMetadataCache) InvalidateObject(ctx context.Context, bucketID, objectKey string) error {
	key := objectCacheKey(bucketID, objectKey)
	return c.client.Del(ctx, key).Err()
}

func (c *redisMetadataCache) GetBucket(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error) {
	key := bucketCacheKey(workspaceID, bucketName)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var bucket model.Bucket
	if err := json.Unmarshal([]byte(val), &bucket); err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (c *redisMetadataCache) SetBucket(ctx context.Context, bucket *model.Bucket, ttl time.Duration) error {
	key := bucketCacheKey(bucket.WorkspaceID.String(), bucket.Name)
	data, err := json.Marshal(bucket)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}
