package provider

import (
	"context"
)

// AIProvider represents the standard abstraction for all AI inference and embedding models.
// Implementations (e.g., openai, anthropic, mock) must satisfy this interface.
type AIProvider interface {
	ID() string
	Name() string
	Health(ctx context.Context) error

	// GenerateText generates a single completion response for a given prompt.
	GenerateText(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)

	// ChatCompletion supports multi-turn conversational models.
	ChatCompletion(ctx context.Context, req ChatRequest) (*GenerateResponse, error)

	// GenerateEmbeddings converts text into a high-dimensional vector space.
	GenerateEmbeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error)

	// Future capabilities (Tool Calling, Structured Output) can be added here without breaking existing integrations.
}
