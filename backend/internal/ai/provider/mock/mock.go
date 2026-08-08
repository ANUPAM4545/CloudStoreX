package mock

import (
	"context"
	"time"

	"github.com/cloudstorex/backend/internal/ai/provider"
)

type MockProvider struct {
	id   string
	name string
}

func NewMockProvider(id, name string) provider.AIProvider {
	return &MockProvider{
		id:   id,
		name: name,
	}
}

func (m *MockProvider) ID() string {
	return m.id
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Health(ctx context.Context) error {
	return nil
}

func (m *MockProvider) GenerateText(ctx context.Context, req provider.GenerateRequest) (*provider.GenerateResponse, error) {
	start := time.Now()
	
	// Deterministic mock logic
	text := "Mock generated text based on: " + req.Prompt

	return &provider.GenerateResponse{
		Text:         text,
		InputTokens:  10,
		OutputTokens: 20,
		TotalTokens:  30,
		Cost:         0.001,
		Model:        req.Model,
		ProviderName: m.name,
		Latency:      time.Since(start),
	}, nil
}

func (m *MockProvider) ChatCompletion(ctx context.Context, req provider.ChatRequest) (*provider.GenerateResponse, error) {
	start := time.Now()

	text := "Mock chat completion response."
	if len(req.Messages) > 0 {
		text = "Mock response to: " + req.Messages[len(req.Messages)-1].Content
	}

	return &provider.GenerateResponse{
		Text:         text,
		InputTokens:  15,
		OutputTokens: 25,
		TotalTokens:  40,
		Cost:         0.002,
		Model:        req.Model,
		ProviderName: m.name,
		Latency:      time.Since(start),
	}, nil
}

func (m *MockProvider) GenerateEmbeddings(ctx context.Context, req provider.EmbeddingRequest) (*provider.EmbeddingResponse, error) {
	start := time.Now()
	
	// Return a deterministic dummy vector (e.g., 1536 dims for openai, let's just do 3 for mock)
	embedding := []float32{0.1, 0.2, 0.3}

	return &provider.EmbeddingResponse{
		Embedding:    embedding,
		InputTokens:  10,
		TotalTokens:  10,
		Cost:         0.0001,
		Model:        req.Model,
		ProviderName: m.name,
		Latency:      time.Since(start),
	}, nil
}
