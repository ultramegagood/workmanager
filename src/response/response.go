package response

import "app/src/model"

type SuccessWithProject struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Project model.Project `json:"project"`
}

type SuccessWithTask struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Task    model.Task `json:"task"`
}

type SuccessWithGroup struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Section   model.Section `json:"section"`
}

type SuccessWithPaginateProjects struct {
	Code         int             `json:"code"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	Results      []model.Project `json:"results"`
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	TotalPages   int64           `json:"total_pages"`
	TotalResults int64           `json:"total_results"`
}

type SuccessWithPaginateTasks struct {
	Code         int          `json:"code"`
	Status       string       `json:"status"`
	Message      string       `json:"message"`
	Results      []model.Task `json:"results"`
	Page         int          `json:"page"`
	Limit        int          `json:"limit"`
	TotalPages   int64        `json:"total_pages"`
	TotalResults int64        `json:"total_results"`
}

type SuccessWithPaginateGroups struct {
	Code         int             `json:"code"`
	Status       string          `json:"status"`
	Message      string          `json:"message"`
	Results      []model.Section `json:"results"`
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	TotalPages   int64           `json:"total_pages"`
	TotalResults int64           `json:"total_results"`
}

type Common struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
type SuccessWithCurrentUser struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
}

type SuccessWithUser struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
}

type SuccessWithTokens struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
	Tokens  Tokens     `json:"tokens"`
}

type SuccessWithData[T any] struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
type SuccessWithPaginate[T any] struct {
	Code         int    `json:"code"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Results      []T    `json:"results"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	TotalPages   int64  `json:"total_pages"`
	TotalResults int64  `json:"total_results"`
}

type ErrorDetails struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
