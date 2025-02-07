package model

import (
	"github.com/google/uuid"
)

type TaskProject struct {
	TaskID   uuid.UUID `gorm:"not null" json:"task_id"`
	ProjectID uuid.UUID `gorm:"not null" json:"project_id"`
	GroupID  uuid.UUID `gorm:"not null" json:"group_id"`
}

