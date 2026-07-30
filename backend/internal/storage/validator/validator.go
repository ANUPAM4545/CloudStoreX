package validator

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

var (
	bucketRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$`)
)

// ValidateBucketName checks whether a bucket name adheres to naming conventions (similar to S3/MinIO rules).
func ValidateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("bucket name must be between 3 and 63 characters")
	}
	if !bucketRegex.MatchString(name) {
		return errors.New("bucket name must contain only lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen")
	}
	if net.ParseIP(name) != nil {
		return errors.New("bucket name cannot be formatted as an IP address")
	}
	return nil
}

// ValidateObjectKey checks whether an object key is valid.
func ValidateObjectKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("object key cannot be empty")
	}
	if len(key) > 1024 {
		return errors.New("object key exceeds maximum length of 1024 characters")
	}
	if strings.Contains(key, "\x00") {
		return errors.New("object key contains invalid null characters")
	}
	if strings.Contains(key, "..") {
		return errors.New("object key cannot contain path traversal sequence '..'")
	}
	return nil
}

// ValidateUpload checks whether upload parameters (key, size, maxSize) are valid.
func ValidateUpload(key string, size, maxSize int64) error {
	if err := ValidateObjectKey(key); err != nil {
		return err
	}
	if size > maxSize && maxSize > 0 {
		return errors.New("file size exceeds maximum allowed upload size")
	}
	return nil
}
