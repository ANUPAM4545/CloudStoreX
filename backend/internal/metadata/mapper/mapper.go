package mapper

import (
	"github.com/cloudstorex/backend/internal/metadata/dto"
	"github.com/cloudstorex/backend/internal/metadata/model"
)

// ToObjectMetaDTO converts a database model to an API DTO.
func ToObjectMetaDTO(obj *model.Object) dto.ObjectMetaDTO {
	tags := make(map[string]string)
	for _, t := range obj.Tags {
		tags[t.Key] = t.Value
	}

	meta := make(map[string]string)
	for _, m := range obj.Metadata {
		meta[m.Key] = m.Value
	}

	version := ""
	if obj.VersionID != nil {
		version = *obj.VersionID
	}

	owner := ""
	if obj.OwnerID != nil {
		owner = obj.OwnerID.String()
	}

	return dto.ObjectMetaDTO{
		ID:                obj.ID.String(),
		BucketID:          obj.BucketID.String(),
		ObjectKey:         obj.ObjectKey,
		ProviderObjectKey: obj.ProviderObjectKey,
		SizeBytes:         obj.SizeBytes,
		MimeType:          obj.MimeType,
		ETag:              obj.ETag,
		ProviderID:        obj.ProviderID,
		VersionID:         version,
		OwnerID:           owner,
		StorageClass:      obj.StorageClass,
		Status:            obj.Status,
		Tags:              tags,
		Metadata:          meta,
		CreatedAt:         obj.CreatedAt,
		UpdatedAt:         obj.UpdatedAt,
	}
}

// ToObjectMetaDTOList converts a list of database models to API DTOs.
func ToObjectMetaDTOList(objs []model.Object) []dto.ObjectMetaDTO {
	list := make([]dto.ObjectMetaDTO, len(objs))
	for i, obj := range objs {
		list[i] = ToObjectMetaDTO(&obj)
	}
	return list
}
