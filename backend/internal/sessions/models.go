package sessions

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionActive  SessionStatus = "active"
	SessionRevoked SessionStatus = "revoked"
	SessionExpired SessionStatus = "expired"
)

type Session struct {
	ID        uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID     `gorm:"type:uuid;not null;index"`
	TokenHash string        `gorm:"type:varchar(255);not null;uniqueIndex"`
	IPAddress string        `gorm:"type:varchar(45)"`
	UserAgent string        `gorm:"type:varchar(255)"`
	DeviceID  string        `gorm:"type:varchar(255)"`
	Status    SessionStatus `gorm:"type:varchar(50);not null"`
	ExpiresAt time.Time     `gorm:"not null"`
	LastSeen  time.Time     `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeviceHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	DeviceID  string    `gorm:"type:varchar(255);not null"`
	UserAgent string    `gorm:"type:varchar(255)"`
	IPAddress string    `gorm:"type:varchar(45)"`
	FirstSeen time.Time
	LastSeen  time.Time
}
