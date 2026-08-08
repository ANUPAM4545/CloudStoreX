package forecasting

import (
	"context"

	"github.com/cloudstorex/backend/internal/ai/provider"
)

type ForecastService interface {
	PredictCapacity(ctx context.Context, workspaceID string) (string, error)
	PredictReliability(ctx context.Context, providerName string) (string, error)
}

type forecastService struct {
	aiManager provider.Manager
}

func NewForecastService(m provider.Manager) ForecastService {
	return &forecastService{
		aiManager: m,
	}
}

func (s *forecastService) PredictCapacity(ctx context.Context, workspaceID string) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	req := provider.GenerateRequest{
		Prompt:      "Forecast storage capacity exhaustion for workspace " + workspaceID,
		Model:       "default-forecasting-model",
		Temperature: 0.1,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

func (s *forecastService) PredictReliability(ctx context.Context, providerName string) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	req := provider.GenerateRequest{
		Prompt:      "Predict reliability degradation and failover probability for provider " + providerName,
		Model:       "default-forecasting-model",
		Temperature: 0.1,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
