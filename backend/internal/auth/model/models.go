package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderType string

const (
	ProviderLocal    ProviderType = "local"
	ProviderGoogle   ProviderType = "google"
	ProviderGitHub   ProviderType = "github"
	ProviderEntraID  ProviderType = "entraid"
	ProviderOIDC     ProviderType = "oidc"
	ProviderSAML     ProviderType = "saml"
)

// Identity represents a user's linked identity from an external provider (or local)
type Identity struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index"`
	Provider       ProviderType   `gorm:"type:varchar(50);not null;index:idx_provider_subject,unique"`
	ProviderSubject string         `gorm:"type:varchar(255);not null;index:idx_provider_subject,unique"` // e.g. the sub claim or external user ID
	Email          string         `gorm:"type:varchar(255)"`
	RawProfile     []byte         `gorm:"type:jsonb"` // JSON representation of the profile received
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type AuthRequest struct {
	State        string `json:"state"`
	Nonce        string `json:"nonce"`
	Provider     ProviderType `json:"provider"`
	RedirectURI  string `json:"redirect_uri"`
	CreatedAt    time.Time
}

type AuthResult struct {
	UserID       uuid.UUID
	Email        string
	Provider     ProviderType
	Subject      string
	IsNewUser    bool
	AccessToken  string
	RefreshToken string
}
