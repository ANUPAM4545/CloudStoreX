package mfa

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MFAType string

const (
	MFATypeTOTP MFAType = "totp"
	MFATypeSMS  MFAType = "sms"
)

type MFAConfig struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Type      MFAType        `gorm:"type:varchar(50);not null"`
	Secret    string         `gorm:"type:varchar(255);not null"`
	Verified  bool           `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RecoveryCode struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CodeHash  string    `gorm:"type:varchar(255);not null"`
	Used      bool      `gorm:"default:false"`
	UsedAt    *time.Time
	CreatedAt time.Time
}
