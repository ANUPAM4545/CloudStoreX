package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LifecycleAction string

const (
	ActionExpire       LifecycleAction = "Expire"
	ActionArchive      LifecycleAction = "Archive"
	ActionTransition   LifecycleAction = "Transition"
	ActionRestore      LifecycleAction = "Restore"
	ActionDeleteMarker LifecycleAction = "DeleteMarker"
)

// LifecycleRule defines automated actions on objects based on age and prefix.
type LifecycleRule struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	BucketID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Prefix      string         `gorm:"type:varchar(255);not null;index"`
	Action      LifecycleAction `gorm:"type:varchar(50);not null"`
	AgeDays     int            `gorm:"not null"`
	Status      string         `gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
