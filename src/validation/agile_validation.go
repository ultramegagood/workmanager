package validation

import "github.com/google/uuid"

type CreateProject struct {
	Title    string      `json:"title" validate:"required,max=50" example:"fake name"`
	GroupIDs []uuid.UUID `json:"group_ids"` // Было []string
}

type CreateGroup struct {
	Title     string      `json:"title" validate:"required,max=50" example:"fake name"`
	TaskIDs   []uuid.UUID `json:"task_ids"`  // Было []string
	ProjectID uuid.UUID   `json:"project_id" validate:"required,uuid"`
}

type CreateTask struct {
	Title        string      `json:"title" validate:"required,max=50" example:"fake task"`
	Description  string      `json:"description"`
	ProjectID    uuid.UUID   `json:"project_id" validate:"required,uuid"`  // Добавил теги
	ParentTaskID *uuid.UUID  `json:"parent_task_id,omitempty"` // Если есть parent task
}
