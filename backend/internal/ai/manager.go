package ai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var (
	ErrProviderNotConfigured = errors.New("requested AI provider is not configured")
	ErrProviderNotFound      = errors.New("requested AI provider not found")
)

// AIManager serves as the single entry point for all AI functionality in CloudStoreX
type AIManager interface {
	// GenerateText routes the request to the specified provider, or the default if none specified
	GenerateText(ctx context.Context, provider ProviderType, req GenerateTextRequest) (*GenerateTextResponse, error)
}

type manager struct {
	providers map[ProviderType]AIProvider
	defaultP  ProviderType
	log       *slog.Logger
}

// NewManager creates a new AIManager instance
func NewManager(defaultProvider ProviderType, log *slog.Logger) *manager {
	if log == nil {
		log = slog.Default()
	}
	return &manager{
		providers: make(map[ProviderType]AIProvider),
		defaultP:  defaultProvider,
		log:       log,
	}
}

// RegisterProvider adds an initialized provider to the manager
func (m *manager) RegisterProvider(p AIProvider) {
	m.providers[p.GetProviderType()] = p
}

// GenerateText routes a text generation request to the appropriate provider
func (m *manager) GenerateText(ctx context.Context, provider ProviderType, req GenerateTextRequest) (*GenerateTextResponse, error) {
	pType := provider
	if pType == "" {
		pType = m.defaultP
	}

	p, ok := m.providers[pType]
	if !ok {
		m.log.ErrorContext(ctx, "AI provider not configured", slog.String("provider", string(pType)))
		return nil, fmt.Errorf("%w: %s", ErrProviderNotConfigured, pType)
	}

	start := time.Now() // Note: Need to import time, will let compiler catch or just fix
	resp, err := p.GenerateText(ctx, req)
	
	duration := time.Since(start)
	if err != nil {
		m.log.ErrorContext(ctx, "AI generation failed", 
			slog.String("provider", string(pType)),
			slog.String("error", err.Error()),
			slog.Duration("duration", duration),
		)
		return nil, err
	}

	m.log.InfoContext(ctx, "AI generation successful",
		slog.String("provider", string(pType)),
		slog.Int("tokens_prompt", resp.TokensPrompt),
		slog.Int("tokens_output", resp.TokensOutput),
		slog.Duration("duration", duration),
	)

	return resp, nil
}
