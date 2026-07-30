package storage

import (
	"errors"
	"fmt"
)

var (
	ErrProviderNotFound   = errors.New("storage provider not found")
	ErrBucketNotFound     = errors.New("bucket not found")
	ErrObjectNotFound     = errors.New("object not found")
	ErrInvalidProvider    = errors.New("invalid storage provider configuration")
	ErrUploadFailed       = errors.New("object upload failed")
	ErrDownloadFailed     = errors.New("object download failed")
	ErrDeleteFailed       = errors.New("object delete failed")
	ErrBucketExists       = errors.New("bucket already exists")
	ErrInvalidRequest     = errors.New("invalid storage request")
	ErrMultipartFailed    = errors.New("multipart upload operation failed")
	ErrUnsupportedFeature = errors.New("provider does not support requested feature")
)

// DomainError wraps domain errors with additional operational context.
type DomainError struct {
	Op       string
	Bucket   string
	Key      string
	Provider string
	Err      error
}

func (e *DomainError) Error() string {
	if e.Key != "" {
		return fmt.Sprintf("storage error [%s] provider=%s bucket=%s key=%s: %v", e.Op, e.Provider, e.Bucket, e.Key, e.Err)
	}
	if e.Bucket != "" {
		return fmt.Sprintf("storage error [%s] provider=%s bucket=%s: %v", e.Op, e.Provider, e.Bucket, e.Err)
	}
	return fmt.Sprintf("storage error [%s] provider=%s: %v", e.Op, e.Provider, e.Err)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError is a helper to wrap storage errors with context.
func NewDomainError(op, provider, bucket, key string, err error) error {
	return &DomainError{
		Op:       op,
		Provider: provider,
		Bucket:   bucket,
		Key:      key,
		Err:      err,
	}
}
