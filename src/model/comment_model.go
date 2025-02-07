package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	ID            uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	TaskID        uuid.UUID `gorm:"not null" json:"task_id"`
	Body          string    `gorm:"not null" json:"body"`
	CitateID      *uuid.UUID `json:"citate_id,omitempty"`
	ReplyToID     *uuid.UUID `json:"reply_to_id,omitempty"`
	IsEdited      bool      `gorm:"default:false" json:"is_edited"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}
func (token *Comment) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}