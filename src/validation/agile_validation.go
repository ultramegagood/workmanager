package validation

import "github.com/google/uuid"

type CreateProject struct {
	Title string `json:"title" validate:"required,max=50" example:"fake name"`
}

type CreateGroup struct {
	Title     string      `json:"title" validate:"required,max=50" example:"fake name"`
	TaskIDs   []uuid.UUID `json:"task_ids"` // Было []string
	ProjectID uuid.UUID   `json:"project_id" validate:"required,uuid"`
}

type CreateTask struct {
	Title        string     `json:"title" validate:"required,max=50" example:"fake task"`
	Description  string     `json:"description"`
	AssignedTo   *uuid.UUID `json:"assigned_to"`
	ProjectID    uuid.UUID  `json:"project_id" validate:"required,uuid"` // Добавил теги
	ParentTaskID *uuid.UUID `json:"parent_task_id,omitempty"`            // Если есть parent task
}
type CreateComment struct {
	Body   string    `json:"body" validate:"required" example:"fake comment"`
	TaskID uuid.UUID `json:"task_id" validate:"required,uuid"` // Добавил теги
}

type CreateUserGroup struct {
	TeamTitle string    `json:"team_title" validate:"required,max=100" example:"Developers"`
	UserID    uuid.UUID `json:"user_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type AddUserToGroup struct {
	GroupID uuid.UUID `json:"user_group_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID  uuid.UUID `json:"user_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type ReassignTaskValidation struct {
	TaskID uuid.UUID `json:"task_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	NewUserID  uuid.UUID `json:"new_user_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type AddGroupToProject struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	GroupID   uuid.UUID `json:"group_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type AddGroupToTask struct {
	TaskID  uuid.UUID `json:"task_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	GroupID uuid.UUID `json:"group_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type GetUsersInGroup struct {
	GroupID uuid.UUID `json:"group_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdateTaskTitleOrDescription struct {
	TaskID      uuid.UUID `json:"task_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title" validate:"required" example:"Title task"`
	Description string    `json:"description" validate:"required" example:"Lorem ipsum"`
}
