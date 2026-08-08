package apikeys

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, scopes []string, expiresAt *time.Time) (*APIKey, string, error)
	RevokeAPIKey(ctx context.Context, id uuid.UUID) error
	ValidateAPIKey(ctx context.Context, prefix, secret string) (*APIKey, error)
}

type defaultService struct {
	// Repo goes here
}

func NewService() Service {
	return &defaultService{}
}

func (s *defaultService) CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, scopes []string, expiresAt *time.Time) (*APIKey, string, error) {
	// Generate prefix (e.g. csx_abc123) and secret
	prefixBytes := make([]byte, 8)
	secretBytes := make([]byte, 32)
	rand.Read(prefixBytes)
	rand.Read(secretBytes)
	
	prefix := "csx_" + base64.RawURLEncoding.EncodeToString(prefixBytes)
	secret := base64.RawURLEncoding.EncodeToString(secretBytes)
	
	key := &APIKey{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Prefix:    prefix,
		KeyHash:   "mock-argon2-hash-" + secret, // Mocking hash
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	}
	
	// Full key to return to user once
	fullKey := prefix + "." + secret
	
	// Save `key` to DB here...
	
	return key, fullKey, nil
}

func (s *defaultService) RevokeAPIKey(ctx context.Context, id uuid.UUID) error {
	// Mark as revoked in DB
	return nil
}

func (s *defaultService) ValidateAPIKey(ctx context.Context, prefix, secret string) (*APIKey, error) {
	// Fetch by prefix, compare hash with secret
	if prefix == "" || secret == "" {
		return nil, errors.New("invalid key format")
	}
	return &APIKey{UserID: uuid.New()}, nil
}
