package router

import (
	"app/src/config"
	"app/src/service"
	"app/src/validation"
	"github.com/gofiber/fiber/v2"	
	"gorm.io/gorm"
	"github.com/gofiber/contrib/websocket"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	tokenService := service.NewTokenService(db, validate, userService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	taskService := service.NewTaskService(db, validate)
	v1 := app.Group("/v1")
	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, tokenService, emailService)
	ProjectRoutes(v1, taskService,userService)
	UserRoutes(v1, userService, tokenService, taskService)
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/tasks", websocket.New(func(c *websocket.Conn) {
		taskService.HandleTaskUpdates(c)
	}))
	if !config.IsProd {
		DocsRoutes(v1)
	}
}
