package model

import (
	"time"

	"github.com/google/uuid"
)


type TaskHistory struct {
	TaskID    uuid.UUID `gorm:"not null" json:"task_id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	Action    string    `gorm:"not null" json:"action"`
	Timestamp time.Time `gorm:"autoCreateTime:milli" json:"timestamp"`
}