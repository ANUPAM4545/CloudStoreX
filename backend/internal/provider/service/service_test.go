package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/provider/dto"
	"github.com/cloudstorex/backend/internal/provider/model"
	"github.com/cloudstorex/backend/internal/provider/service"
	"github.com/google/uuid"
)

// Mock Repository
type mockProviderRepo struct {
	providers map[string]*model.Provider
	defaultID string
}

func (m *mockProviderRepo) Create(ctx context.Context, p *model.Provider) error {
	m.providers[p.ID.String()] = p
	return nil
}
func (m *mockProviderRepo) GetByID(ctx context.Context, id string) (*model.Provider, error) {
	if p, ok := m.providers[id]; ok {
		return p, nil
	}
	return nil, nil // Return err in reality
}
func (m *mockProviderRepo) GetByName(ctx context.Context, wsID, name string) (*model.Provider, error) {
	return nil, nil
}
func (m *mockProviderRepo) List(ctx context.Context, wsID string) ([]*model.Provider, error) {
	var res []*model.Provider
	for _, p := range m.providers {
		if p.WorkspaceID.String() == wsID {
			res = append(res, p)
		}
	}
	return res, nil
}
func (m *mockProviderRepo) GetDefault(ctx context.Context, wsID string) (*model.Provider, error) {
	if m.defaultID != "" {
		return m.providers[m.defaultID], nil
	}
	return nil, nil
}
func (m *mockProviderRepo) Update(ctx context.Context, p *model.Provider) error {
	m.providers[p.ID.String()] = p
	return nil
}
func (m *mockProviderRepo) UpdateHealth(ctx context.Context, id string, status model.ProviderStatus, health model.ProviderHealth, latency int64, errStr string) error {
	if p, ok := m.providers[id]; ok {
		p.Status = status
		p.Health = health
		p.LatencyMs = latency
		p.LastError = errStr
		now := time.Now()
		p.LastHealthCheck = &now
	}
	return nil
}
func (m *mockProviderRepo) SetDefault(ctx context.Context, wsID, id string) error {
	m.defaultID = id
	return nil
}
func (m *mockProviderRepo) Enable(ctx context.Context, id string) error {
	if p, ok := m.providers[id]; ok {
		p.IsEnabled = true
	}
	return nil
}
func (m *mockProviderRepo) Disable(ctx context.Context, id string) error {
	if p, ok := m.providers[id]; ok {
		p.IsEnabled = false
	}
	return nil
}
func (m *mockProviderRepo) Delete(ctx context.Context, id string) error {
	delete(m.providers, id)
	return nil
}

// Mock Validator
type mockValidator struct {
	latency int64
	err     error
}
func (m *mockValidator) ValidateConnection(ctx context.Context, p *model.Provider, extended bool) (int64, error) {
	return m.latency, m.err
}

func TestProviderService_CreateAndList(t *testing.T) {
	repo := &mockProviderRepo{providers: make(map[string]*model.Provider)}
	val := &mockValidator{latency: 10, err: nil}
	svc := service.NewProviderService(repo, val)
	
	wsID := uuid.New()
	
	req := dto.CreateProviderRequest{
		WorkspaceID:  wsID,
		ProviderName: "test-aws",
		ProviderType: "AWS_S3",
		IsDefault:    true,
	}
	
	p, err := svc.CreateProvider(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.ProviderName != "test-aws" {
		t.Errorf("expected test-aws, got %s", p.ProviderName)
	}
	
	providers, err := svc.ListProviders(context.Background(), wsID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(providers) != 1 {
		t.Errorf("expected 1 provider, got %d", len(providers))
	}
}

func TestProviderService_Validation(t *testing.T) {
	repo := &mockProviderRepo{providers: make(map[string]*model.Provider)}
	val := &mockValidator{latency: 15, err: nil}
	svc := service.NewProviderService(repo, val)
	
	p := &model.Provider{
		ID: uuid.New(),
		WorkspaceID: uuid.New(),
		Status: model.StatusPending,
		Health: model.HealthUnknown,
	}
	repo.providers[p.ID.String()] = p
	
	err := svc.ValidateProviderConnection(context.Background(), p.ID.String(), false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	updatedP := repo.providers[p.ID.String()]
	if updatedP.Status != model.StatusReady {
		t.Errorf("expected READY status, got %s", updatedP.Status)
	}
	if updatedP.Health != model.HealthHealthy {
		t.Errorf("expected HEALTHY health, got %s", updatedP.Health)
	}
}
