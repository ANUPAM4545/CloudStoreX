package repository

import (
	"context"

	"github.com/cloudstorex/backend/internal/iam/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Organizations
	CreateOrganization(ctx context.Context, org *model.Organization) error
	GetOrganization(ctx context.Context, id uuid.UUID) (*model.Organization, error)
	ListOrganizations(ctx context.Context) ([]model.Organization, error)
	
	// Departments
	CreateDepartment(ctx context.Context, dept *model.Department) error
	GetDepartment(ctx context.Context, id uuid.UUID) (*model.Department, error)
	ListDepartments(ctx context.Context, orgID uuid.UUID) ([]model.Department, error)

	// Teams
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeam(ctx context.Context, id uuid.UUID) (*model.Team, error)
	ListTeams(ctx context.Context, orgID uuid.UUID) ([]model.Team, error)

	// UserGroups
	CreateUserGroup(ctx context.Context, group *model.UserGroup) error
	GetUserGroup(ctx context.Context, id uuid.UUID) (*model.UserGroup, error)
	ListUserGroups(ctx context.Context, orgID uuid.UUID) ([]model.UserGroup, error)

	// Memberships
	AddMembership(ctx context.Context, mem *model.Membership) error
	RemoveMembership(ctx context.Context, id uuid.UUID) error
	ListMemberships(ctx context.Context, entityType model.EntityType, entityID uuid.UUID) ([]model.Membership, error)
	ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]model.Membership, error)

	// Invitations
	CreateInvitation(ctx context.Context, inv *model.Invitation) error
	GetInvitation(ctx context.Context, id uuid.UUID) (*model.Invitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*model.Invitation, error)
	UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status model.InvitationStatus) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// Organizations
func (r *gormRepository) CreateOrganization(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *gormRepository) GetOrganization(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *gormRepository) ListOrganizations(ctx context.Context) ([]model.Organization, error) {
	var orgs []model.Organization
	if err := r.db.WithContext(ctx).Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

// Departments
func (r *gormRepository) CreateDepartment(ctx context.Context, dept *model.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *gormRepository) GetDepartment(ctx context.Context, id uuid.UUID) (*model.Department, error) {
	var dept model.Department
	if err := r.db.WithContext(ctx).First(&dept, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *gormRepository) ListDepartments(ctx context.Context, orgID uuid.UUID) ([]model.Department, error) {
	var depts []model.Department
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

// Teams
func (r *gormRepository) CreateTeam(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *gormRepository) GetTeam(ctx context.Context, id uuid.UUID) (*model.Team, error) {
	var team model.Team
	if err := r.db.WithContext(ctx).First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *gormRepository) ListTeams(ctx context.Context, orgID uuid.UUID) ([]model.Team, error) {
	var teams []model.Team
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

// UserGroups
func (r *gormRepository) CreateUserGroup(ctx context.Context, group *model.UserGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *gormRepository) GetUserGroup(ctx context.Context, id uuid.UUID) (*model.UserGroup, error) {
	var group model.UserGroup
	if err := r.db.WithContext(ctx).First(&group, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *gormRepository) ListUserGroups(ctx context.Context, orgID uuid.UUID) ([]model.UserGroup, error) {
	var groups []model.UserGroup
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// Memberships
func (r *gormRepository) AddMembership(ctx context.Context, mem *model.Membership) error {
	return r.db.WithContext(ctx).Create(mem).Error
}

func (r *gormRepository) RemoveMembership(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Membership{}, "id = ?", id).Error
}

func (r *gormRepository) ListMemberships(ctx context.Context, entityType model.EntityType, entityID uuid.UUID) ([]model.Membership, error) {
	var mems []model.Membership
	if err := r.db.WithContext(ctx).Where("entity_type = ? AND entity_id = ?", entityType, entityID).Find(&mems).Error; err != nil {
		return nil, err
	}
	return mems, nil
}

func (r *gormRepository) ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]model.Membership, error) {
	var mems []model.Membership
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&mems).Error; err != nil {
		return nil, err
	}
	return mems, nil
}

// Invitations
func (r *gormRepository) CreateInvitation(ctx context.Context, inv *model.Invitation) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *gormRepository) GetInvitation(ctx context.Context, id uuid.UUID) (*model.Invitation, error) {
	var inv model.Invitation
	if err := r.db.WithContext(ctx).First(&inv, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *gormRepository) GetInvitationByToken(ctx context.Context, token string) (*model.Invitation, error) {
	var inv model.Invitation
	if err := r.db.WithContext(ctx).First(&inv, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *gormRepository) UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status model.InvitationStatus) error {
	return r.db.WithContext(ctx).Model(&model.Invitation{}).Where("id = ?", id).Update("status", status).Error
}
