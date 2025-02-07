package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	TaskID     uuid.UUID `gorm:"not null" json:"task_id"`
	Body       string    `json:"body"`
	ActionType string    `gorm:"not null" json:"action_type"`
	EntityType string    `gorm:"not null" json:"entity_type"`
	EntityID   uuid.UUID `gorm:"not null" json:"entity_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}
func (token *AuditLog) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}