package recommendations

import (
	"context"

	"github.com/cloudstorex/backend/internal/ai/provider"
)

type OptimizationService interface {
	OptimizeCosts(ctx context.Context, workspaceID string) (string, error)
	AdvisePolicies(ctx context.Context, workspaceID string) (string, error)
}

type optimizationService struct {
	aiManager provider.Manager
}

func NewOptimizationService(m provider.Manager) OptimizationService {
	return &optimizationService{
		aiManager: m,
	}
}

func (s *optimizationService) OptimizeCosts(ctx context.Context, workspaceID string) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	// Mocking passing workspace telemetry
	req := provider.GenerateRequest{
		Prompt:      "Analyze cost efficiency for workspace " + workspaceID,
		Model:       "default-recommendation-model",
		Temperature: 0.2,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

func (s *optimizationService) AdvisePolicies(ctx context.Context, workspaceID string) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	// Mocking analyzing routing rules
	req := provider.GenerateRequest{
		Prompt:      "Analyze routing policies and suggest improvements for workspace " + workspaceID,
		Model:       "default-recommendation-model",
		Temperature: 0.2,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
