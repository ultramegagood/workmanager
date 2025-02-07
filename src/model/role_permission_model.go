package model

import "github.com/google/uuid"



type RolePermission struct {
	RoleID uuid.UUID `gorm:"not null" json:"role_id"`
	UserID uuid.UUID `gorm:"not null" json:"user_id"`
}
