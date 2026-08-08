package security

import (
	"context"

	"github.com/cloudstorex/backend/internal/ai/provider"
)

type IntelligenceService interface {
	AnalyzeAnomalies(ctx context.Context, orgID string) (string, error)
}

type intelligenceService struct {
	aiManager provider.Manager
}

func NewIntelligenceService(m provider.Manager) IntelligenceService {
	return &intelligenceService{
		aiManager: m,
	}
}

func (s *intelligenceService) AnalyzeAnomalies(ctx context.Context, orgID string) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	req := provider.GenerateRequest{
		Prompt:      "Analyze security audit logs for anomalies, login spikes, or abuse in organization " + orgID,
		Model:       "default-security-model",
		Temperature: 0.1,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
