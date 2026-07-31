package service

import (
	"context"
	"encoding/json"

	"github.com/cloudstorex/backend/internal/provider/dto"
	"github.com/cloudstorex/backend/internal/provider/model"
	"github.com/cloudstorex/backend/internal/provider/repository"
	"github.com/google/uuid"
)

// ProviderService manages provider persistence, health, and configuration.
type ProviderService interface {
	CreateProvider(ctx context.Context, req dto.CreateProviderRequest) (*dto.ProviderDTO, error)
	UpdateProvider(ctx context.Context, id string, req dto.UpdateProviderRequest) (*dto.ProviderDTO, error)
	EnableProvider(ctx context.Context, id string) error
	DisableProvider(ctx context.Context, id string) error
	GetProvider(ctx context.Context, id string) (*dto.ProviderDTO, error)
	ListProviders(ctx context.Context, workspaceID string) ([]*dto.ProviderDTO, error)
	GetDefaultProvider(ctx context.Context, workspaceID string) (*dto.ProviderDTO, error)
	GetDefaultProviderID(ctx context.Context, workspaceID string) (string, error)
	SetDefaultProvider(ctx context.Context, workspaceID, id string) error
	ReportHealth(ctx context.Context, id string, status model.ProviderStatus, health model.ProviderHealth, latencyMs int64, lastError string) error
	ValidateProviderConnection(ctx context.Context, id string, extended bool) error
	RefreshProviderHealth(ctx context.Context, id string) error
	DeleteProvider(ctx context.Context, id string) error
}

type defaultProviderService struct {
	repo      repository.Repository
	validator ProviderValidator
}

func NewProviderService(repo repository.Repository, validator ProviderValidator) ProviderService {
	return &defaultProviderService{repo: repo, validator: validator}
}

func (s *defaultProviderService) CreateProvider(ctx context.Context, req dto.CreateProviderRequest) (*dto.ProviderDTO, error) {
	capsBytes, err := json.Marshal(req.Capabilities)
	if err != nil {
		return nil, err
	}

	p := &model.Provider{
		ID:               uuid.New(),
		WorkspaceID:      req.WorkspaceID,
		ProviderName:     req.ProviderName,
		ProviderType:     req.ProviderType,
		Endpoint:         req.Endpoint,
		Region:           req.Region,
		BucketPrefix:     req.BucketPrefix,
		IsDefault:        req.IsDefault,
		IsEnabled:        true,
		Status:           model.StatusPending,
		Health:           model.HealthUnknown,
		CredentialRef:    req.CredentialRef,
		CapabilitiesJSON: string(capsBytes),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	if req.IsDefault {
		_ = s.repo.SetDefault(ctx, p.WorkspaceID.String(), p.ID.String())
	}

	return s.mapToDTO(p), nil
}

func (s *defaultProviderService) GetProvider(ctx context.Context, id string) (*dto.ProviderDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapToDTO(p), nil
}

func (s *defaultProviderService) UpdateProvider(ctx context.Context, id string, req dto.UpdateProviderRequest) (*dto.ProviderDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Endpoint != nil {
		p.Endpoint = *req.Endpoint
	}
	if req.Region != nil {
		p.Region = *req.Region
	}
	if req.BucketPrefix != nil {
		p.BucketPrefix = *req.BucketPrefix
	}
	if req.CredentialRef != nil {
		p.CredentialRef = *req.CredentialRef
	}
	if req.Capabilities != nil {
		capsBytes, _ := json.Marshal(req.Capabilities)
		p.CapabilitiesJSON = string(capsBytes)
	}

	p.Status = model.StatusPending
	p.Health = model.HealthUnknown
	
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}

	return s.mapToDTO(p), nil
}

func (s *defaultProviderService) EnableProvider(ctx context.Context, id string) error {
	return s.repo.Enable(ctx, id)
}

func (s *defaultProviderService) DisableProvider(ctx context.Context, id string) error {
	return s.repo.Disable(ctx, id)
}

func (s *defaultProviderService) ListProviders(ctx context.Context, workspaceID string) ([]*dto.ProviderDTO, error) {
	providers, err := s.repo.List(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.ProviderDTO, len(providers))
	for i, p := range providers {
		dtos[i] = s.mapToDTO(p)
	}
	return dtos, nil
}

func (s *defaultProviderService) GetDefaultProvider(ctx context.Context, workspaceID string) (*dto.ProviderDTO, error) {
	p, err := s.repo.GetDefault(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return s.mapToDTO(p), nil
}

func (s *defaultProviderService) GetDefaultProviderID(ctx context.Context, workspaceID string) (string, error) {
	p, err := s.repo.GetDefault(ctx, workspaceID)
	if err != nil {
		return "", err
	}
	return p.ProviderName, nil // ProviderName is currently used as ID in the runtime registry
}

func (s *defaultProviderService) SetDefaultProvider(ctx context.Context, workspaceID, id string) error {
	return s.repo.SetDefault(ctx, workspaceID, id)
}

func (s *defaultProviderService) ReportHealth(ctx context.Context, id string, status model.ProviderStatus, health model.ProviderHealth, latencyMs int64, lastError string) error {
	return s.repo.UpdateHealth(ctx, id, status, health, latencyMs, lastError)
}

func (s *defaultProviderService) ValidateProviderConnection(ctx context.Context, id string, extended bool) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Update status to Validating
	_ = s.repo.UpdateHealth(ctx, id, model.StatusValidating, p.Health, p.LatencyMs, p.LastError)

	latency, err := s.validator.ValidateConnection(ctx, p, extended)
	if err != nil {
		_ = s.repo.UpdateHealth(ctx, id, model.StatusFailed, model.HealthUnavailable, latency, err.Error())
		return err
	}

	// Update status to Ready
	_ = s.repo.UpdateHealth(ctx, id, model.StatusReady, model.HealthHealthy, latency, "")
	return nil
}

func (s *defaultProviderService) RefreshProviderHealth(ctx context.Context, id string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	latency, err := s.validator.ValidateConnection(ctx, p, false)
	if err != nil {
		_ = s.repo.UpdateHealth(ctx, id, p.Status, model.HealthUnavailable, latency, err.Error())
		return err
	}

	_ = s.repo.UpdateHealth(ctx, id, p.Status, model.HealthHealthy, latency, "")
	return nil
}

func (s *defaultProviderService) DeleteProvider(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *defaultProviderService) mapToDTO(p *model.Provider) *dto.ProviderDTO {
	var caps dto.ProviderCapabilities
	_ = json.Unmarshal([]byte(p.CapabilitiesJSON), &caps)

	return &dto.ProviderDTO{
		ID:           p.ID,
		WorkspaceID:  p.WorkspaceID,
		ProviderName: p.ProviderName,
		ProviderType: p.ProviderType,
		Endpoint:     p.Endpoint,
		Region:       p.Region,
		BucketPrefix: p.BucketPrefix,
		IsDefault:    p.IsDefault,
		IsEnabled:    p.IsEnabled,
		Status:       string(p.Status),
		Health:       string(p.Health),
		LastHealthCheck: p.LastHealthCheck,
		LatencyMs:    p.LatencyMs,
		LastError:    p.LastError,
		Capabilities: caps,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
