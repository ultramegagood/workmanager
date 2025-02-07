package model

import "github.com/google/uuid"


type UserGroup struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	TeamTitle string    `json:"team_title"`
}