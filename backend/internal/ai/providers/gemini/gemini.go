package gemini

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/ai"
)

type provider struct {
	apiKey string
}

// NewProvider creates a new Gemini provider wrapper
func NewProvider(ctx context.Context, apiKey string) (ai.AIProvider, error) {
	return &provider{
		apiKey: apiKey,
	}, nil
}

func (p *provider) GetProviderType() ai.ProviderType {
	return ai.ProviderGemini
}

func (p *provider) GenerateText(ctx context.Context, req ai.GenerateTextRequest) (*ai.GenerateTextResponse, error) {
	return nil, errors.New("gemini provider not fully implemented yet")
}
