package model

import (
	"github.com/google/uuid"
)

type Project struct {
	ID       uuid.UUID   `gorm:"primaryKey;not null" json:"id"`
	Title    string      `gorm: json:"title" `
	GroupIDs []uuid.UUID `gorm:"type:uuid[]" json:"group_ids"`
}
