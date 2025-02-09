package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID            uuid.UUID  `gorm:"primaryKey;" json:"id"`
	ProjectID     uuid.UUID  `gorm:"" json:"project_id"`
	Title         string     `gorm:"" json:"title"`
	Description   string     `gorm:"" json:"description"`
	UserGroup     uuid.UUID  `json:"user_group"`
	Status        string     `gorm:"" json:"status"`
	Priority      string     `json:"priority"`
	DueDate       time.Time  `json:"due_date"`
	CreatedAt     time.Time  `gorm:"autoCreateTime:milli" json:"created_at"`
	AssignedTo    uuid.UUID  `json:"assigned_to"`
	ParentTaskID  *uuid.UUID `json:"parent_task_id,omitempty"`
	EstimatedTime int        `json:"estimated_time"`
	SpentTime     int        `json:"spent_time"`
}

func (token *Task) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}
