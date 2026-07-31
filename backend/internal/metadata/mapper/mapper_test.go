package mapper_test

import (
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/metadata/mapper"
	"github.com/cloudstorex/backend/internal/metadata/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToObjectMetaDTO(t *testing.T) {
	objID := uuid.New()
	bucketID := uuid.New()
	ownerID := uuid.New()
	versionID := "v1.0"
	now := time.Now()

	obj := &model.Object{
		ID:                objID,
		BucketID:          bucketID,
		ObjectKey:         "test.txt",
		ProviderObjectKey: "prov/test.txt",
		SizeBytes:         1024,
		MimeType:          "text/plain",
		ETag:              "abcdef",
		ProviderID:        "minio",
		VersionID:         &versionID,
		OwnerID:           &ownerID,
		StorageClass:      "STANDARD",
		Status:            "ACTIVE",
		CreatedAt:         now,
		UpdatedAt:         now,
		Tags: []model.ObjectTag{
			{Key: "env", Value: "prod"},
		},
		Metadata: []model.ObjectMetadata{
			{Key: "custom", Value: "val"},
		},
	}

	dtoObj := mapper.ToObjectMetaDTO(obj)

	assert.Equal(t, objID.String(), dtoObj.ID)
	assert.Equal(t, bucketID.String(), dtoObj.BucketID)
	assert.Equal(t, "test.txt", dtoObj.ObjectKey)
	assert.Equal(t, 1024, int(dtoObj.SizeBytes))
	assert.Equal(t, "v1.0", dtoObj.VersionID)
	assert.Equal(t, ownerID.String(), dtoObj.OwnerID)
	assert.Equal(t, "prod", dtoObj.Tags["env"])
	assert.Equal(t, "val", dtoObj.Metadata["custom"])
}

func TestToObjectMetaDTOList(t *testing.T) {
	objs := []model.Object{
		{ID: uuid.New(), ObjectKey: "1.txt"},
		{ID: uuid.New(), ObjectKey: "2.txt"},
	}

	dtos := mapper.ToObjectMetaDTOList(objs)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "1.txt", dtos[0].ObjectKey)
	assert.Equal(t, "2.txt", dtos[1].ObjectKey)
}
