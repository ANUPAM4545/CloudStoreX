package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/cloudstorex/backend/internal/iam/model"
	"github.com/cloudstorex/backend/internal/iam/repository"
	"github.com/google/uuid"
)

type Service interface {
	CreateOrganization(ctx context.Context, name, domain string) (*model.Organization, error)
	CreateDepartment(ctx context.Context, orgID uuid.UUID, name string) (*model.Department, error)
	CreateTeam(ctx context.Context, orgID uuid.UUID, deptID *uuid.UUID, name string) (*model.Team, error)
	CreateUserGroup(ctx context.Context, orgID uuid.UUID, name string) (*model.UserGroup, error)
	
	AddMembership(ctx context.Context, userID uuid.UUID, entityType model.EntityType, entityID uuid.UUID, role string) (*model.Membership, error)
	
	InviteUser(ctx context.Context, orgID uuid.UUID, email, role string, ttl time.Duration) (*model.Invitation, error)
	AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error
}

type defaultService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &defaultService{repo: repo}
}

func (s *defaultService) CreateOrganization(ctx context.Context, name, domain string) (*model.Organization, error) {
	if name == "" {
		return nil, errors.New("organization name is required")
	}
	org := &model.Organization{
		ID:     uuid.New(),
		Name:   name,
		Domain: domain,
	}
	if err := s.repo.CreateOrganization(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *defaultService) CreateDepartment(ctx context.Context, orgID uuid.UUID, name string) (*model.Department, error) {
	dept := &model.Department{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
	}
	if err := s.repo.CreateDepartment(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *defaultService) CreateTeam(ctx context.Context, orgID uuid.UUID, deptID *uuid.UUID, name string) (*model.Team, error) {
	team := &model.Team{
		ID:             uuid.New(),
		OrganizationID: orgID,
		DepartmentID:   deptID,
		Name:           name,
	}
	if err := s.repo.CreateTeam(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *defaultService) CreateUserGroup(ctx context.Context, orgID uuid.UUID, name string) (*model.UserGroup, error) {
	group := &model.UserGroup{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
	}
	if err := s.repo.CreateUserGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *defaultService) AddMembership(ctx context.Context, userID uuid.UUID, entityType model.EntityType, entityID uuid.UUID, role string) (*model.Membership, error) {
	mem := &model.Membership{
		ID:         uuid.New(),
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Role:       role,
	}
	if err := s.repo.AddMembership(ctx, mem); err != nil {
		return nil, err
	}
	return mem, nil
}

func (s *defaultService) InviteUser(ctx context.Context, orgID uuid.UUID, email, role string, ttl time.Duration) (*model.Invitation, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	inv := &model.Invitation{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		Token:          token,
		ExpiresAt:      time.Now().Add(ttl),
		Status:         model.InvitationPending,
	}
	
	if err := s.repo.CreateInvitation(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *defaultService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error {
	inv, err := s.repo.GetInvitationByToken(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}
	if inv.Status != model.InvitationPending {
		return errors.New("invitation is not pending")
	}
	if time.Now().After(inv.ExpiresAt) {
		_ = s.repo.UpdateInvitationStatus(ctx, inv.ID, model.InvitationExpired)
		return errors.New("invitation expired")
	}

	mem := &model.Membership{
		ID:         uuid.New(),
		UserID:     userID,
		EntityType: model.EntityOrganization,
		EntityID:   inv.OrganizationID,
		Role:       inv.Role,
	}
	if err := s.repo.AddMembership(ctx, mem); err != nil {
		return err
	}

	return s.repo.UpdateInvitationStatus(ctx, inv.ID, model.InvitationAccepted)
}
