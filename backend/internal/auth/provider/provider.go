package provider

import (
	"context"
	"net/http"

	"github.com/cloudstorex/backend/internal/auth/model"
)

type Provider interface {
	Type() model.ProviderType
	
	// GetLoginURL returns the provider-specific redirect URL for starting the auth flow
	GetLoginURL(ctx context.Context, state, nonce string) (string, error)
	
	// HandleCallback exchanges the code for a token and retrieves the user profile
	HandleCallback(ctx context.Context, req *http.Request) (*model.AuthResult, error)
}
