package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudstorex/backend/internal/auth/model"
	"github.com/google/uuid"
)

// MockProvider provides a lightweight mock for standalone execution of enterprise SSO flows
type MockProvider struct {
	providerType model.ProviderType
}

func NewMockProvider(ptype model.ProviderType) Provider {
	return &MockProvider{providerType: ptype}
}

func (m *MockProvider) Type() model.ProviderType {
	return m.providerType
}

func (m *MockProvider) GetLoginURL(ctx context.Context, state, nonce string) (string, error) {
	// In a real system, this would build a URL to Google, Entra, etc.
	return fmt.Sprintf("https://mock-idp.example.com/auth?state=%s&nonce=%s&provider=%s", state, nonce, m.providerType), nil
}

func (m *MockProvider) HandleCallback(ctx context.Context, req *http.Request) (*model.AuthResult, error) {
	// Mock parsing the authorization code and returning a dummy identity
	code := req.URL.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("missing authorization code")
	}

	return &model.AuthResult{
		UserID:       uuid.Nil, // to be populated by the service
		Email:        fmt.Sprintf("user-%s@example.com", code),
		Provider:     m.providerType,
		Subject:      fmt.Sprintf("sub-%s", code),
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
	}, nil
}
