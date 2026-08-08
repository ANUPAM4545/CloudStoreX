package mfa

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GenerateSecret(ctx context.Context, userID uuid.UUID, email string) (string, string, error)
	VerifyTOTP(ctx context.Context, userID uuid.UUID, code string) (bool, error)
	ValidateRecoveryCode(ctx context.Context, userID uuid.UUID, code string) (bool, error)
}

type defaultService struct {
	// Normally we would have an MFA repository here to store the config and recovery codes
}

func NewService() Service {
	return &defaultService{}
}

func (s *defaultService) GenerateSecret(ctx context.Context, userID uuid.UUID, email string) (string, string, error) {
	// Generate a 20-byte random secret
	secretBytes := make([]byte, 20)
	if _, err := rand.Read(secretBytes); err != nil {
		return "", "", err
	}
	secret := base32.StdEncoding.EncodeToString(secretBytes)
	
	// Create provisioning URI (for QR Code generation on the frontend)
	uri := fmt.Sprintf("otpauth://totp/CloudStoreX:%s?secret=%s&issuer=CloudStoreX", email, secret)
	
	// In reality, we'd save this to the DB as unverified.
	return secret, uri, nil
}

func (s *defaultService) VerifyTOTP(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	// In reality, this would fetch the user's secret from the DB and validate the 6-digit TOTP
	if code == "" || len(code) != 6 {
		return false, errors.New("invalid code format")
	}
	// Mock successful validation
	return true, nil
}

func (s *defaultService) ValidateRecoveryCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	// In reality, this would check the DB for unused recovery codes matching the hash, and mark as used.
	if code == "" {
		return false, errors.New("empty recovery code")
	}
	return true, nil
}
