package identity

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service interface {
	Register(ctx context.Context, email, password, firstName, lastName string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type service struct {
	repo         Repository
	tokenService *TokenService
}

func NewService(repo Repository, tokenService *TokenService) Service {
	return &service{repo: repo, tokenService: tokenService}
}

func (s *service) Register(ctx context.Context, email, password, firstName, lastName string) (*User, error) {
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	valid, err := VerifyPassword(password, user.PasswordHash)
	if err != nil || !valid {
		return "", ErrInvalidCredentials
	}

	// Generate JWT
	token, err := s.tokenService.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}
