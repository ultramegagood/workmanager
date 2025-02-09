package validation

type CreateProject struct {
	Title    string   `json:"title" validate:"required,max=50" example:"fake name"`
	GroupIDs []string `json:"group_ids"`
}
type CreateGroup struct {
	Title   string   `json:"title" validate:"required,max=50" example:"fake name"`
	TaskIDs []string `json:"task_ids"`
}
type CreateTask struct {
	Title       string   `json:"title" validate:"required,max=50" example:"fake name"`
	Description string   `json:"description"`
	TaskIDs     []string `json:"task_ids"`
	ProjectID 	string
}
