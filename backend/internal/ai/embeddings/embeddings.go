package embeddings

import (
	"context"

	"github.com/cloudstorex/backend/internal/ai/provider"
)

type Store interface {
	Index(ctx context.Context, id string, text string, embedding []float32) error
	Search(ctx context.Context, queryEmbedding []float32, topK int) ([]SearchResult, error)
}

type SearchResult struct {
	ID    string
	Score float32
}

type Service interface {
	GenerateAndIndex(ctx context.Context, id string, text string) error
	Search(ctx context.Context, query string, topK int) ([]SearchResult, error)
}

type service struct {
	aiProvider provider.AIProvider
	store      Store
}

func NewService(p provider.AIProvider, store Store) Service {
	return &service{
		aiProvider: p,
		store:      store,
	}
}

func (s *service) GenerateAndIndex(ctx context.Context, id string, text string) error {
	req := provider.EmbeddingRequest{
		Input: text,
		Model: "default-embedding-model",
	}
	resp, err := s.aiProvider.GenerateEmbeddings(ctx, req)
	if err != nil {
		return err
	}
	return s.store.Index(ctx, id, text, resp.Embedding)
}

func (s *service) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	req := provider.EmbeddingRequest{
		Input: query,
		Model: "default-embedding-model",
	}
	resp, err := s.aiProvider.GenerateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}
	return s.store.Search(ctx, resp.Embedding, topK)
}
