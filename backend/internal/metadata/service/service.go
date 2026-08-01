package service

import (
	"context"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/events"
	"github.com/cloudstorex/backend/internal/metadata/model"
	"github.com/cloudstorex/backend/internal/metadata/repository"
	"github.com/google/uuid"
)

type MetadataService interface {
	CreateObjectMetadata(ctx context.Context, workspaceID string, bucketID uuid.UUID, objectKey, providerObjectKey string, size int64, mimeType, etag, providerID string, ownerID *uuid.UUID, tags map[string]string, meta map[string]string) (*model.Object, error)
	UpdateObjectMetadata(ctx context.Context, id string, mimeType string, storageClass string, tags map[string]string, custom map[string]string) (*model.Object, error)
	FindObjectByKey(ctx context.Context, bucketID, objectKey string) (*model.Object, error)
	FindObjectByID(ctx context.Context, id string) (*model.Object, error)
	GetObjectMetadata(ctx context.Context, id string) (map[string]string, error)
	SearchObjects(ctx context.Context, query dto.SearchQuery) ([]dto.ObjectMetaDTO, int64, error)
	FindObjectsForExpiration(ctx context.Context, bucketID string, prefix string, olderThan time.Time) ([]model.Object, error)
	TagObject(ctx context.Context, id string, tags map[string]string) error
	UntagObject(ctx context.Context, id string, keys []string) error
	SoftDeleteObject(ctx context.Context, workspaceID, id string) error
	RestoreObject(ctx context.Context, id string) error
	ListObjectVersions(ctx context.Context, id string) ([]dto.ObjectVersionDTO, error)
	
	CreateBucket(ctx context.Context, workspaceID uuid.UUID, providerID, name, region string) (*model.Bucket, error)
	FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error)
	ListBuckets(ctx context.Context, workspaceID string) ([]model.Bucket, error)
}

type metadataService struct {
	repo      repository.MetadataRepository
	publisher events.Publisher
}

func NewMetadataService(repo repository.MetadataRepository, publisher events.Publisher) MetadataService {
	return &metadataService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *metadataService) CreateObjectMetadata(ctx context.Context, workspaceID string, bucketID uuid.UUID, objectKey, providerObjectKey string, size int64, mimeType, etag, providerID string, ownerID *uuid.UUID, tags map[string]string, meta map[string]string) (*model.Object, error) {
	// 1. Check if object already exists
	existingObj, err := s.repo.FindObjectByKey(ctx, bucketID.String(), objectKey)
	
	var objID uuid.UUID
	var nextVersionNum int = 1
	var obj *model.Object

	if err != nil {
		// Object doesn't exist, create it
		objID = uuid.New()
		
		var objTags []model.ObjectTag
		for k, v := range tags {
			objTags = append(objTags, model.ObjectTag{ID: uuid.New(), ObjectID: objID, Key: k, Value: v})
		}
		
		var objMeta []model.ObjectMetadata
		for k, v := range meta {
			objMeta = append(objMeta, model.ObjectMetadata{ID: uuid.New(), ObjectID: objID, Key: k, Value: v})
		}

		obj = &model.Object{
			ID:                objID,
			BucketID:          bucketID,
			ObjectKey:         objectKey,
			ProviderObjectKey: providerObjectKey,
			SizeBytes:         size,
			MimeType:          mimeType,
			ETag:              etag,
			ProviderID:        providerID,
			OwnerID:           ownerID,
			Status:            "ACTIVE",
			Tags:              objTags,
			Metadata:          objMeta,
		}

		if err := s.repo.CreateObject(ctx, obj); err != nil {
			return nil, err
		}
	} else {
		// Object exists. Update it and increment version.
		objID = existingObj.ID
		// In a real app we might want to query max version, but for simplicity we can just count versions or fetch max.
		// For now we'll just set it to len(existingObj.Versions) + 1 if preloaded, but FindObjectByKey might not preload.
		// Let's assume VersionNumber is generated or we just use Unix timestamp for VersionNumber if we can't reliably get max.
		// Better: we can query the max version in the repo.
		// For now, let's just use time.Now().UnixNano() as a naive version number, or just 1. Let's use Unix()
		nextVersionNum = int(time.Now().Unix())
		
		existingObj.SizeBytes = size
		existingObj.MimeType = mimeType
		existingObj.ETag = etag
		existingObj.ProviderObjectKey = providerObjectKey
		existingObj.ProviderID = providerID
		existingObj.UpdatedAt = time.Now()
		
		if err := s.repo.UpdateObject(ctx, existingObj); err != nil {
			return nil, err
		}
		obj = existingObj
	}

	// 2. Create the Object Version
	version := &model.ObjectVersion{
		ID:                uuid.New(),
		ObjectID:          objID,
		VersionNumber:     nextVersionNum,
		ProviderID:        providerID,
		ProviderObjectKey: providerObjectKey,
		SizeBytes:         size,
		ETag:              etag,
		IsCurrent:         true,
	}

	if err := s.repo.CreateObjectVersion(ctx, version); err != nil {
		return nil, err
	}

	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventObjectCreated,
		Timestamp: time.Now(),
		BucketID:  bucketID.String(),
		ObjectID:  objID.String(),
		Payload: map[string]interface{}{
			"size_bytes":   size,
			"workspace_id": workspaceID,
		},
	})

	return obj, nil
}

func (s *metadataService) FindObjectByKey(ctx context.Context, bucketID, objectKey string) (*model.Object, error) {
	return s.repo.FindObjectByKey(ctx, bucketID, objectKey)
}

func (s *metadataService) SearchObjects(ctx context.Context, query dto.SearchQuery) ([]dto.ObjectMetaDTO, int64, error) {
	objs, total, err := s.repo.SearchObjects(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// For now, mapping inside the service layer because we didn't inject the mapper explicitly
	// Assuming mapper is used in HTTP handlers, but we'll return DTOs here to avoid leaking GORM structs in complex APIs
	// Actually, returning DTOs from Service is best practice.
	var results []dto.ObjectMetaDTO
	for _, obj := range objs {
		t := make(map[string]string)
		for _, tag := range obj.Tags {
			t[tag.Key] = tag.Value
		}
		m := make(map[string]string)
		for _, meta := range obj.Metadata {
			m[meta.Key] = meta.Value
		}
		
		results = append(results, dto.ObjectMetaDTO{
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
		})
	}
	
	return results, total, nil
}

func (s *metadataService) FindObjectsForExpiration(ctx context.Context, bucketID string, prefix string, olderThan time.Time) ([]model.Object, error) {
	return s.repo.FindObjectsForExpiration(ctx, bucketID, prefix, olderThan)
}

func (s *metadataService) SoftDeleteObject(ctx context.Context, workspaceID, id string) error {
	// Let's fetch the object size before deleting to include it in the event
	obj, err := s.repo.FindObjectByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteObject(ctx, id); err != nil {
		return err
	}

	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventObjectDeleted,
		Timestamp: time.Now(),
		ObjectID:  id,
		Payload: map[string]interface{}{
			"size_bytes":   obj.SizeBytes,
			"workspace_id": workspaceID,
		},
	})

	return nil
}

func (s *metadataService) RestoreObject(ctx context.Context, id string) error {
	obj, err := s.repo.FindObjectByID(ctx, id)
	if err != nil {
		return err
	}

	obj.IsDeleted = false
	obj.TrashTimestamp = nil
	obj.DeletedAt = nil
	if err := s.repo.UpdateObject(ctx, obj); err != nil {
		return err
	}

	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventObjectRestored,
		Timestamp: time.Now(),
		ObjectID:  id,
	})
	return nil
}

func (s *metadataService) ListObjectVersions(ctx context.Context, id string) ([]dto.ObjectVersionDTO, error) {
	versions, err := s.repo.ListObjectVersions(ctx, id)
	if err != nil {
		return nil, err
	}
	var dtos []dto.ObjectVersionDTO
	for _, v := range versions {
		dtos = append(dtos, dto.ObjectVersionDTO{
			ID:                v.ID.String(),
			ObjectID:          v.ObjectID.String(),
			VersionNumber:     v.VersionNumber,
			ProviderID:        v.ProviderID,
			ProviderObjectKey: v.ProviderObjectKey,
			SizeBytes:         v.SizeBytes,
			ETag:              v.ETag,
			IsCurrent:         v.IsCurrent,
			CreatedAt:         v.CreatedAt,
		})
	}
	return dtos, nil
}

func (s *metadataService) FindObjectByID(ctx context.Context, id string) (*model.Object, error) {
	return s.repo.FindObjectByID(ctx, id)
}

func (s *metadataService) UpdateObjectMetadata(ctx context.Context, id string, mimeType string, storageClass string, tags map[string]string, custom map[string]string) (*model.Object, error) {
	obj, err := s.repo.FindObjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if mimeType != "" {
		obj.MimeType = mimeType
	}
	if storageClass != "" {
		obj.StorageClass = storageClass
	}
	
	// Update custom metadata (replace completely or merge? Usually replace or merge. Let's replace for simplicity as per standard API unless merge is specified)
	// We should probably rely on the Repository to clear old ones in a transaction, but for now we'll do an in-memory update.
	// Actually, GORM association replacement is safer.
	if custom != nil {
		var newMeta []model.ObjectMetadata
		for k, v := range custom {
			newMeta = append(newMeta, model.ObjectMetadata{ID: uuid.New(), ObjectID: obj.ID, Key: k, Value: v})
		}
		if err := s.repo.DB().WithContext(ctx).Model(obj).Association("Metadata").Replace(newMeta); err != nil {
			return nil, err
		}
		obj.Metadata = newMeta
	}
	
	if tags != nil {
		var newTags []model.ObjectTag
		for k, v := range tags {
			newTags = append(newTags, model.ObjectTag{ID: uuid.New(), ObjectID: obj.ID, Key: k, Value: v})
		}
		if err := s.repo.DB().WithContext(ctx).Model(obj).Association("Tags").Replace(newTags); err != nil {
			return nil, err
		}
		obj.Tags = newTags
	}

	if err := s.repo.UpdateObject(ctx, obj); err != nil {
		return nil, err
	}

	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventObjectUpdated,
		Timestamp: time.Now(),
		ObjectID:  id,
	})

	return obj, nil
}

func (s *metadataService) GetObjectMetadata(ctx context.Context, id string) (map[string]string, error) {
	obj, err := s.repo.FindObjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, meta := range obj.Metadata {
		m[meta.Key] = meta.Value
	}
	return m, nil
}

func (s *metadataService) TagObject(ctx context.Context, id string, tags map[string]string) error {
	obj, err := s.repo.FindObjectByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Append/Update tags
	tagMap := make(map[string]string)
	for _, t := range obj.Tags {
		tagMap[t.Key] = t.Value
	}
	for k, v := range tags {
		tagMap[k] = v
	}
	
	var newTags []model.ObjectTag
	for k, v := range tagMap {
		newTags = append(newTags, model.ObjectTag{ID: uuid.New(), ObjectID: obj.ID, Key: k, Value: v})
	}
	
	if err := s.repo.DB().WithContext(ctx).Model(obj).Association("Tags").Replace(newTags); err != nil {
		return err
	}

	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventObjectTagged,
		Timestamp: time.Now(),
		ObjectID:  id,
	})
	return nil
}

func (s *metadataService) UntagObject(ctx context.Context, id string, keys []string) error {
	obj, err := s.repo.FindObjectByID(ctx, id)
	if err != nil {
		return err
	}
	
	keyToRemove := make(map[string]bool)
	for _, k := range keys {
		keyToRemove[k] = true
	}
	
	var newTags []model.ObjectTag
	for _, t := range obj.Tags {
		if !keyToRemove[t.Key] {
			newTags = append(newTags, t)
		}
	}
	
	if err := s.repo.DB().WithContext(ctx).Model(obj).Association("Tags").Replace(newTags); err != nil {
		return err
	}
	
	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventObjectUntagged,
		Timestamp: time.Now(),
		ObjectID:  id,
	})
	return nil
}

func (s *metadataService) CreateBucket(ctx context.Context, workspaceID uuid.UUID, providerID, name, region string) (*model.Bucket, error) {
	bucket := &model.Bucket{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		ProviderID:  providerID,
		Name:        name,
		Region:      region,
	}

	if err := s.repo.CreateBucket(ctx, bucket); err != nil {
		return nil, err
	}

	s.publisher.Publish(ctx, events.Event{
		ID:        uuid.New().String(),
		Type:      events.EventBucketCreated,
		Timestamp: time.Now(),
		BucketID:  bucket.ID.String(),
	})

	return bucket, nil
}

func (s *metadataService) FindBucketByName(ctx context.Context, workspaceID, bucketName string) (*model.Bucket, error) {
	return s.repo.FindBucketByName(ctx, workspaceID, bucketName)
}

func (s *metadataService) ListBuckets(ctx context.Context, workspaceID string) ([]model.Bucket, error) {
	return s.repo.ListBuckets(ctx, workspaceID)
}
