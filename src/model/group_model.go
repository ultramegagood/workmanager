package model

import "github.com/google/uuid"
type Group struct {
	ID       uuid.UUID   `gorm:"primaryKey;not null" json:"id"`
	Title    string      `gorm:"not null" json:"title"`
	TaskIDs  []uuid.UUID `gorm:"type:uuid[]" json:"task_ids"`
	ProjectID uuid.UUID  `gorm:"not null" json:"project_id"`
}
