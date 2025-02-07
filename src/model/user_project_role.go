package model

import "github.com/google/uuid"

type UserProjectRole struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	ProjectID uuid.UUID `gorm:"not null" json:"project_id"`
}