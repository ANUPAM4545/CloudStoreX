package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudstorex/backend/internal/auth/model"
	"github.com/cloudstorex/backend/internal/auth/provider"
	"github.com/google/uuid"
)

type Service interface {
	RegisterProvider(p provider.Provider)
	GetLoginURL(ctx context.Context, ptype model.ProviderType, redirectURI string) (string, error)
	HandleCallback(ctx context.Context, ptype model.ProviderType, req *http.Request) (*model.AuthResult, error)
}

type authService struct {
	providers map[model.ProviderType]provider.Provider
}

func NewService() Service {
	return &authService{
		providers: make(map[model.ProviderType]provider.Provider),
	}
}

func (s *authService) RegisterProvider(p provider.Provider) {
	s.providers[p.Type()] = p
}

func (s *authService) GetLoginURL(ctx context.Context, ptype model.ProviderType, redirectURI string) (string, error) {
	p, ok := s.providers[ptype]
	if !ok {
		return "", fmt.Errorf("unsupported provider: %s", ptype)
	}

	state := uuid.NewString()
	nonce := uuid.NewString()

	// In a real implementation, we would store state & nonce in a session/cache to verify during callback
	return p.GetLoginURL(ctx, state, nonce)
}

func (s *authService) HandleCallback(ctx context.Context, ptype model.ProviderType, req *http.Request) (*model.AuthResult, error) {
	p, ok := s.providers[ptype]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", ptype)
	}

	// 1. Let the provider handle the token exchange and profile fetching
	res, err := p.HandleCallback(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. Here we would link the identity to a user in the database (or create a new user)
	// (Omitted for brevity in this mock layer, would use Identity Repository)

	return res, nil
}
