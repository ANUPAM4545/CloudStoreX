package anthropic

import (
	"context"
	"errors"

	"github.com/cloudstorex/backend/internal/ai"
)

type provider struct {
	apiKey string
}

// NewProvider creates a new Anthropic provider wrapper
func NewProvider(apiKey string) ai.AIProvider {
	return &provider{
		apiKey: apiKey,
	}
}

func (p *provider) GetProviderType() ai.ProviderType {
	return ai.ProviderAnthropic
}

func (p *provider) GenerateText(ctx context.Context, req ai.GenerateTextRequest) (*ai.GenerateTextResponse, error) {
	return nil, errors.New("anthropic provider not fully implemented yet")
}
