package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func ProjectRoutes(v1 fiber.Router, t service.TaskService, u service.UserService) {
	taskController := controller.NewTaskController(t)


	project := v1.Group("/projects")

	// Проекты
	project.Post("/", m.Auth(u), taskController.CreateProject)
	// Задачи
	project.Post("/:projectId/tasks", m.Auth(u), taskController.CreateTask)
	// Группы
	project.Post("/:projectId/groups", m.Auth(u), taskController.CreateGroup)
}
