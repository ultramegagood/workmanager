package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type UserGroup struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	TeamTitle string    `json:"team_title"`
}
func (user *UserGroup) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
