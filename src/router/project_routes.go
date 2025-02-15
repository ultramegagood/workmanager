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
	// Задачи
	v1.Post("/create-task", m.Auth(u), taskController.CreateTask)
	// Группы
	v1.Post("/group", m.Auth(u), taskController.CreateGroup)
	v1.Get("/tasks/:taskID/users", m.Auth(u), taskController.GetUsersWithAccess)
	v1.Post("/user-groups", m.Auth(u), taskController.CreateUserGroup)
	v1.Post("/user-groups/add-user", m.Auth(u), taskController.AddUserToGroup)
	v1.Post("/projects/add-group", m.Auth(u), taskController.AddGroupToProject)
	v1.Post("/tasks/add-group", m.Auth(u), taskController.AddGroupToTask)
	v1.Get("/user-groups", m.Auth(u), taskController.GetUserGroups)
	v1.Post("/user-groups/users", m.Auth(u), taskController.GetUsersInGroup)
}
