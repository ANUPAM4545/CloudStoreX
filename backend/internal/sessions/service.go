package sessions

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateSession(ctx context.Context, userID uuid.UUID, ip, userAgent, deviceID string) (*Session, string, error)
	ValidateSession(ctx context.Context, token string) (*Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
}

type defaultService struct {
	// Repo goes here
}

func NewService() Service {
	return &defaultService{}
}

func (s *defaultService) CreateSession(ctx context.Context, userID uuid.UUID, ip, userAgent, deviceID string) (*Session, string, error) {
	// Generate random token
	tokenBytes := make([]byte, 64)
	rand.Read(tokenBytes)
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	session := &Session{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "mock-hash-" + token,
		IPAddress: ip,
		UserAgent: userAgent,
		DeviceID:  deviceID,
		Status:    SessionActive,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		LastSeen:  time.Now(),
	}

	// Persist session to DB...

	return session, token, nil
}

func (s *defaultService) ValidateSession(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, errors.New("missing session token")
	}
	
	// Mock validation
	// Hash token and fetch from DB
	// Check expiration and status
	return &Session{ID: uuid.New(), UserID: uuid.New(), Status: SessionActive}, nil
}

func (s *defaultService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	// Update DB Status = SessionRevoked
	return nil
}

func (s *defaultService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	// Update DB Status = SessionRevoked where UserID = userID
	return nil
}
