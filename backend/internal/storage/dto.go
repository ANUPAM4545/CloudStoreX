package storage

import "time"

// CreateBucketRequest represents the JSON request body for POST /api/v1/storage/buckets.
type CreateBucketRequest struct {
	Bucket string `json:"bucket" binding:"required"`
}

// BucketDTO represents a storage bucket in API responses.
type BucketDTO struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ObjectDTO represents an object metadata item in API responses.
type ObjectDTO struct {
	Key          string            `json:"key"`
	Bucket       string            `json:"bucket"`
	Size         int64             `json:"size"`
	LastModified time.Time         `json:"last_modified,omitempty"`
	ETag         string            `json:"etag,omitempty"`
	ContentType  string            `json:"content_type,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// UploadResponseDTO represents the successful result of an object upload.
type UploadResponseDTO struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Size   int64  `json:"size"`
	ETag   string `json:"etag,omitempty"`
}

// ToBucketDTO converts an internal Bucket domain model to a BucketDTO.
func ToBucketDTO(b *Bucket) *BucketDTO {
	if b == nil {
		return nil
	}
	return &BucketDTO{
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
	}
}

// ToBucketDTOList converts a slice of internal Buckets to a slice of BucketDTOs.
func ToBucketDTOList(buckets []*Bucket) []*BucketDTO {
	res := make([]*BucketDTO, 0, len(buckets))
	for _, b := range buckets {
		res = append(res, ToBucketDTO(b))
	}
	return res
}

// ToObjectDTO converts an internal Object domain model to an ObjectDTO.
func ToObjectDTO(o *Object) *ObjectDTO {
	if o == nil {
		return nil
	}
	dto := &ObjectDTO{
		Key:          o.Key,
		Bucket:       o.Bucket,
		Size:         o.Size,
		LastModified: o.LastModified,
		ETag:         o.ETag,
	}
	if o.Metadata != nil {
		dto.ContentType = o.Metadata.ContentType
		dto.Metadata = o.Metadata.Custom
	}
	return dto
}

// ToObjectDTOList converts a slice of internal Objects to a slice of ObjectDTOs.
func ToObjectDTOList(objects []*Object) []*ObjectDTO {
	res := make([]*ObjectDTO, 0, len(objects))
	for _, o := range objects {
		res = append(res, ToObjectDTO(o))
	}
	return res
}

// ToUploadResponseDTO converts an internal StorageResponse domain model to an UploadResponseDTO.
func ToUploadResponseDTO(r *StorageResponse) *UploadResponseDTO {
	if r == nil {
		return nil
	}
	return &UploadResponseDTO{
		Bucket: r.Bucket,
		Key:    r.Key,
		Size:   r.Size,
		ETag:   r.ETag,
	}
}
