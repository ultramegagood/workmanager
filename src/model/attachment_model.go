package model

import "github.com/google/uuid"



type Attachment struct {
	ID              uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	TaskID          uuid.UUID `gorm:"not null" json:"task_id"`
	UserID          uuid.UUID `gorm:"not null" json:"user_id"`
	URL             string    `gorm:"not null" json:"url"`
	Type            string    `gorm:"not null" json:"type"`
	Size            int       `json:"size"`
	LinkedTaskID    *uuid.UUID `json:"linked_task_id,omitempty"`
	LinkedCommentID *uuid.UUID `json:"linked_comment_id,omitempty"`
}
