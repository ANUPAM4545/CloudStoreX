package reports

import (
	"context"

	"github.com/cloudstorex/backend/internal/ai/provider"
)

type ReportService interface {
	GenerateExecutiveReport(ctx context.Context, orgID string) (string, error)
	GeneratePlatformReport(ctx context.Context) (string, error)
}

type reportService struct {
	aiManager provider.Manager
}

func NewReportService(m provider.Manager) ReportService {
	return &reportService{
		aiManager: m,
	}
}

func (s *reportService) GenerateExecutiveReport(ctx context.Context, orgID string) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	req := provider.GenerateRequest{
		Prompt:      "Generate an executive summary report for org " + orgID + " covering costs, security, and usage.",
		Model:       "default-report-model",
		Temperature: 0.1,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

func (s *reportService) GeneratePlatformReport(ctx context.Context) (string, error) {
	activeProvider, err := s.aiManager.ActiveProvider()
	if err != nil {
		return "", err
	}

	req := provider.GenerateRequest{
		Prompt:      "Generate a global platform health and operations report.",
		Model:       "default-report-model",
		Temperature: 0.1,
	}

	resp, err := activeProvider.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}
