package router

import (
	"app/src/config"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()
	redisClient := config.RedisClient() // Добавляем Redis

	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	tokenService := service.NewTokenService(db, validate, userService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	taskService := service.NewTaskService(db, validate, redisClient) // Передаём Redis-клиент

	v1 := app.Group("/v1")
	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, tokenService, emailService)
	ProjectRoutes(v1, taskService, userService)
	UserRoutes(v1, userService, tokenService, taskService)

	// Настроим WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/tasks", websocket.New(func(c *websocket.Conn) {
		taskService.HandleTaskUpdates(c)
	}))

	app.Get("/ws/comments", websocket.New(func(c *websocket.Conn) {
		taskService.HandleCommentUpdates(c)
	}))
	app.Get("/ws/projects", websocket.New(func(c *websocket.Conn) {
		taskService.HandleProjectUpdates(c)
	}))
	if !config.IsProd {
		DocsRoutes(v1)
	}
}
