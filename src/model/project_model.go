package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID       uuid.UUID   `gorm:"primaryKey;not null" json:"id"`
	Title    string      `gorm:"not null" json:"title" `
	GroupIDs []uuid.UUID `gorm:"type:uuid[]" json:"group_ids"`
}

func (token *Project) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}