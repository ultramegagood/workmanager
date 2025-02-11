package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID    uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Title string    `json:"title"`
}

func (token *Project) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}
