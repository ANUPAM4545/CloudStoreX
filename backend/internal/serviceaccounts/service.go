package serviceaccounts

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	CreateServiceAccount(ctx context.Context, orgID uuid.UUID, workspaceID *uuid.UUID, name, description string) (*ServiceAccount, string, error)
	RotateSecret(ctx context.Context, id uuid.UUID) (string, error)
	Authenticate(ctx context.Context, clientID, clientSecret string) (*ServiceAccount, error)
}

type defaultService struct {
	// Repo goes here
}

func NewService() Service {
	return &defaultService{}
}

func (s *defaultService) CreateServiceAccount(ctx context.Context, orgID uuid.UUID, workspaceID *uuid.UUID, name, description string) (*ServiceAccount, string, error) {
	clientIDBytes := make([]byte, 16)
	secretBytes := make([]byte, 32)
	rand.Read(clientIDBytes)
	rand.Read(secretBytes)
	
	clientID := "sa_" + hex.EncodeToString(clientIDBytes)
	clientSecret := hex.EncodeToString(secretBytes)
	
	sa := &ServiceAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		WorkspaceID:    workspaceID,
		Name:           name,
		Description:    description,
		ClientID:       clientID,
		ClientSecret:   "mock-hash-" + clientSecret,
		Active:         true,
	}
	
	return sa, clientSecret, nil
}

func (s *defaultService) RotateSecret(ctx context.Context, id uuid.UUID) (string, error) {
	secretBytes := make([]byte, 32)
	rand.Read(secretBytes)
	clientSecret := hex.EncodeToString(secretBytes)
	// Update hash in DB
	return clientSecret, nil
}

func (s *defaultService) Authenticate(ctx context.Context, clientID, clientSecret string) (*ServiceAccount, error) {
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("invalid credentials")
	}
	// Verify from DB...
	return &ServiceAccount{ID: uuid.New(), OrganizationID: uuid.New()}, nil
}
