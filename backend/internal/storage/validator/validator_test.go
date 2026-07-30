package validator

import (
	"testing"
)

func TestValidateBucketName(t *testing.T) {
	validNames := []string{
		"my-bucket",
		"test1234",
		"a-b-c",
		"cloudstorex-default",
	}
	for _, name := range validNames {
		if err := ValidateBucketName(name); err != nil {
			t.Errorf("expected valid bucket name %q, got error: %v", name, err)
		}
	}

	invalidNames := []string{
		"ab",                      // too short
		"-bucket",                 // starts with hyphen
		"bucket-",                 // ends with hyphen
		"MyBucket",                // uppercase
		"192.168.1.1",             // IP address format
		"bucket name with spaces", // spaces
	}
	for _, name := range invalidNames {
		if err := ValidateBucketName(name); err == nil {
			t.Errorf("expected error for invalid bucket name %q, got nil", name)
		}
	}
}

func TestValidateObjectKey(t *testing.T) {
	validKeys := []string{
		"file.txt",
		"folder/subfolder/file.pdf",
		"a",
	}
	for _, key := range validKeys {
		if err := ValidateObjectKey(key); err != nil {
			t.Errorf("expected valid object key %q, got error: %v", key, err)
		}
	}

	invalidKeys := []string{
		"",                  // empty
		"folder/../file.md", // path traversal
		"file\x00name",      // null byte
	}
	for _, key := range invalidKeys {
		if err := ValidateObjectKey(key); err == nil {
			t.Errorf("expected error for invalid object key %q, got nil", key)
		}
	}
}

func TestValidateUpload(t *testing.T) {
	if err := ValidateUpload("valid.txt", 500, 1000); err != nil {
		t.Errorf("expected valid upload, got %v", err)
	}
	if err := ValidateUpload("valid.txt", 2000, 1000); err == nil {
		t.Errorf("expected error for oversized upload, got nil")
	}
}
