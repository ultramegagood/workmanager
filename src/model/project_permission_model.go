package model

import "github.com/google/uuid"


type ProjectPermission struct {
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	ProjectID uuid.UUID `gorm:"not null" json:"project_id"`
	Role      string    `gorm:"not null" json:"role"`
}