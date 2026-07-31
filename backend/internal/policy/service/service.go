package service

import (
	"context"
	"encoding/json"

	"github.com/cloudstorex/backend/internal/policy/dto"
	"github.com/cloudstorex/backend/internal/policy/model"
	"github.com/cloudstorex/backend/internal/policy/repository"
	"gorm.io/datatypes"
)

// PolicyService manages the lifecycle of routing policies.
type PolicyService interface {
	CreatePolicy(ctx context.Context, req dto.CreatePolicyRequest) (*dto.PolicyDTO, error)
	GetPolicy(ctx context.Context, id string) (*dto.PolicyDTO, error)
	ListPolicies(ctx context.Context, workspaceID string) ([]*dto.PolicyDTO, error)
	UpdatePolicy(ctx context.Context, id string, req dto.UpdatePolicyRequest) (*dto.PolicyDTO, error)
	DeletePolicy(ctx context.Context, id string) error
	EnablePolicy(ctx context.Context, id string) error
	DisablePolicy(ctx context.Context, id string) error
	ListRoutingDecisions(ctx context.Context, workspaceID string, limit, offset int) ([]*dto.RoutingDecisionDTO, int64, error)
}

type defaultPolicyService struct {
	repo repository.Repository
}

// NewPolicyService creates a new Policy Service.
func NewPolicyService(repo repository.Repository) PolicyService {
	return &defaultPolicyService{repo: repo}
}

func mapToDTO(p *model.Policy) *dto.PolicyDTO {
	return &dto.PolicyDTO{
		ID:          p.ID.String(),
		WorkspaceID: p.WorkspaceID.String(),
		Name:        p.Name,
		Description: p.Description,
		Priority:    p.Priority,
		Enabled:     p.Enabled,
		RuleType:    string(p.RuleType),
		Conditions:  []byte(p.Conditions),
		Actions:     []byte(p.Actions),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (s *defaultPolicyService) CreatePolicy(ctx context.Context, req dto.CreatePolicyRequest) (*dto.PolicyDTO, error) {
	p := &model.Policy{
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Description: req.Description,
		Priority:    req.Priority,
		Enabled:     true,
		RuleType:    model.RuleType(req.RuleType),
		Conditions:  datatypes.JSON(req.Conditions),
		Actions:     datatypes.JSON(req.Actions),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return mapToDTO(p), nil
}

func (s *defaultPolicyService) GetPolicy(ctx context.Context, id string) (*dto.PolicyDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapToDTO(p), nil
}

func (s *defaultPolicyService) ListPolicies(ctx context.Context, workspaceID string) ([]*dto.PolicyDTO, error) {
	policies, err := s.repo.List(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	
	var res []*dto.PolicyDTO
	for _, p := range policies {
		res = append(res, mapToDTO(p))
	}
	return res, nil
}

func (s *defaultPolicyService) UpdatePolicy(ctx context.Context, id string, req dto.UpdatePolicyRequest) (*dto.PolicyDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Description != "" {
		p.Description = req.Description
	}
	if req.Priority != 0 { // Assume 0 means unchanged, or we could use pointer
		p.Priority = req.Priority
	}
	if req.RuleType != "" {
		p.RuleType = model.RuleType(req.RuleType)
	}
	if req.Conditions != nil {
		p.Conditions = datatypes.JSON(req.Conditions)
	}
	if req.Actions != nil {
		p.Actions = datatypes.JSON(req.Actions)
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return mapToDTO(p), nil
}

func (s *defaultPolicyService) DeletePolicy(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *defaultPolicyService) EnablePolicy(ctx context.Context, id string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	p.Enabled = true
	return s.repo.Update(ctx, p)
}

func (s *defaultPolicyService) DisablePolicy(ctx context.Context, id string) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	p.Enabled = false
	return s.repo.Update(ctx, p)
}

func (s *defaultPolicyService) ListRoutingDecisions(ctx context.Context, workspaceID string, limit, offset int) ([]*dto.RoutingDecisionDTO, int64, error) {
	decisions, total, err := s.repo.ListRoutingDecisions(ctx, workspaceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	
	var res []*dto.RoutingDecisionDTO
	for _, d := range decisions {
		var pid *string
		if d.PolicyID != nil {
			idStr := d.PolicyID.String()
			pid = &idStr
		}
		res = append(res, &dto.RoutingDecisionDTO{
			ID:                 d.ID.String(),
			WorkspaceID:        d.WorkspaceID.String(),
			ObjectKey:          d.ObjectKey,
			ProviderID:         d.ProviderID,
			PolicyID:           pid,
			MatchedRule:        d.MatchedRule,
			IsFallback:         d.IsFallback,
			CandidateProviders: json.RawMessage(d.CandidateProviders),
			DecisionReason:     d.DecisionReason,
			Operation:          d.Operation,
			LatencyMs:          d.LatencyMs,
			Timestamp:          d.Timestamp,
		})
	}
	return res, total, nil
}
