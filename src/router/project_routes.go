package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func ProjectRoutes(v1 fiber.Router, t service.TaskService, u service.UserService) {
	taskController := controller.NewTaskController(t)

	// Проекты
	v1.Post("/projects", m.Auth(u), taskController.CreateProject)
	v1.Get("/projects", m.Auth(u), taskController.GetUserProjects)
	v1.Get("/projects/:projectID/sections", m.Auth(u), taskController.GetSectionsByProject)
	v1.Post("/projects/add-group", m.Auth(u), taskController.AddGroupToProject)

	// Секции
	v1.Post("/projects/section", m.Auth(u), taskController.CreateSection)
	v1.Delete("/sections/:sectionID", m.Auth(u), taskController.DeleteSection)
	v1.Get("/sections", m.Auth(u), taskController.GetSectionsByUser)

	// Задачи
	v1.Post("/tasks", m.Auth(u), taskController.CreateTask)
	v1.Get("/tasks", m.Auth(u), taskController.GetTasks)
	v1.Get("/tasks/:taskID", m.Auth(u), taskController.GetTaskByID)
	v1.Put("/tasks/:taskID", m.Auth(u), taskController.UpdateTaskTitleOrDescription)
	v1.Put("/tasks/:taskID/reassign", m.Auth(u), taskController.ReassignTask)
	v1.Delete("/tasks/:taskID", m.Auth(u), taskController.DeleteTask)
	v1.Get("/tasks/:taskID/users", m.Auth(u), taskController.GetUsersWithAccess)
	v1.Post("/tasks/add-group", m.Auth(u), taskController.AddGroupToTask)

	// Группы пользователей
	v1.Post("/user-groups", m.Auth(u), taskController.CreateUserGroup)
	v1.Post("/user-groups/add-user", m.Auth(u), taskController.AddUserToGroup)
	v1.Get("/user-groups", m.Auth(u), taskController.GetUserGroups)
	v1.Post("/user-groups/users", m.Auth(u), taskController.GetUsersInGroup)

	// Комментарии
	v1.Post("/comments", m.Auth(u), taskController.CommentTask)
}
