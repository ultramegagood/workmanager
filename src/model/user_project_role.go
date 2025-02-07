package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProjectRole struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	ProjectID uuid.UUID `gorm:"not null" json:"project_id"`
}

func (user *UserProjectRole) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
