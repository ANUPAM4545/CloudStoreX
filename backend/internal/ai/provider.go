package ai

import (
	"context"
)

// ProviderType identifies the underlying AI vendor
type ProviderType string

const (
	ProviderOpenAI   ProviderType = "OPENAI"
	ProviderAnthropic ProviderType = "ANTHROPIC"
	ProviderGemini   ProviderType = "GEMINI"
)

// GenerateTextRequest contains the parameters for a text generation request
type GenerateTextRequest struct {
	Prompt      string
	System      string
	Model       string
	Temperature float32
	MaxTokens   int
}

// GenerateTextResponse contains the response from a text generation request
type GenerateTextResponse struct {
	Text         string
	Provider     ProviderType
	ModelUsed    string
	TokensPrompt int
	TokensOutput int
}

// AIProvider defines the contract that all AI vendor SDK wrappers must implement
type AIProvider interface {
	// GenerateText generates a text response based on the prompt
	GenerateText(ctx context.Context, req GenerateTextRequest) (*GenerateTextResponse, error)
	
	// GetProviderType returns the identifier of the provider
	GetProviderType() ProviderType
}
