package audit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SecurityEventType string

const (
	EventLoginSuccess     SecurityEventType = "login_success"
	EventLoginFailure     SecurityEventType = "login_failure"
	EventLogout           SecurityEventType = "logout"
	EventMFAChallenge     SecurityEventType = "mfa_challenge"
	EventMFAEnroll        SecurityEventType = "mfa_enroll"
	EventPermissionChange SecurityEventType = "permission_change"
	EventAPIKeyCreated    SecurityEventType = "api_key_created"
	EventAPIKeyRevoked    SecurityEventType = "api_key_revoked"
	EventRoleAssignment   SecurityEventType = "role_assignment"
	EventSessionRevoked   SecurityEventType = "session_revoked"
)

type SecurityAuditLog struct {
	ID        uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventType SecurityEventType `gorm:"type:varchar(50);not null;index"`
	SubjectID uuid.UUID         `gorm:"type:uuid;index"`
	TargetID  string            `gorm:"type:varchar(255)"` // Optional target entity
	IPAddress string            `gorm:"type:varchar(45)"`
	UserAgent string            `gorm:"type:text"`
	Metadata  []byte            `gorm:"type:jsonb"`
	CreatedAt time.Time         `gorm:"index"`
}

type SecurityAuditLogger interface {
	LogEvent(ctx context.Context, eventType SecurityEventType, subjectID uuid.UUID, targetID, ip, userAgent string, metadata interface{}) error
}

type securityLogger struct {
	db *gorm.DB
}

func NewSecurityLogger(db *gorm.DB) SecurityAuditLogger {
	return &securityLogger{db: db}
}

func (l *securityLogger) LogEvent(ctx context.Context, eventType SecurityEventType, subjectID uuid.UUID, targetID, ip, userAgent string, metadata interface{}) error {
	var metaBytes []byte
	if metadata != nil {
		metaBytes, _ = json.Marshal(metadata)
	}

	log := &SecurityAuditLog{
		ID:        uuid.New(),
		EventType: eventType,
		SubjectID: subjectID,
		TargetID:  targetID,
		IPAddress: ip,
		UserAgent: userAgent,
		Metadata:  metaBytes,
	}

	return l.db.WithContext(ctx).Create(log).Error
}
