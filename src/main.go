package main

import (
	"app/src/config"
	"app/src/database"
	"app/src/middleware"
	"app/src/model"

	"strings"

	"app/src/router"
	"app/src/utils"
	"context"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"

	"gorm.io/gorm"
)

// @title go-fiber-boilerplate API documentation
// @version 1.0.0
// @license.name MIT
// @license.url https://github.com/indrayyana/go-fiber-boilerplate/blob/main/LICENSE
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Example Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := setupFiberApp()
	db := setupDatabase()
	defer closeDatabase(db)
	setupRoutes(app, db)

	address := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)

	// Start server and handle graceful shutdown
	serverErrors := make(chan error, 1)
	go startServer(app, address, serverErrors)
	handleGracefulShutdown(ctx, app, serverErrors)
}

func setupFiberApp() *fiber.App {
	app := fiber.New(config.FiberConfig())

	// Middleware setup
	app.Use("/v1/auth", middleware.LimiterConfig())
	app.Use(middleware.LoggerConfig())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(middleware.RecoverConfig())
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Route().Path, "/ws")
		},
	}))
	app.Use("/ws", func(c *fiber.Ctx) error {
		// Проверяем, является ли запрос WebSocket-подключением
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	router.SetupWebSocketRoutes(app)

	return app
}

func setupDatabase() *gorm.DB {
	db := database.Connect(config.DBHost, config.DBName)
	// Add any additional database setup if needed
	// Выполнение автомиграций

	err := db.AutoMigrate(
		&model.UserSection{}, // Ensure UserSection is included here
		&model.User{},
		&model.Token{},
		&model.Project{},
		&model.Section{},
		&model.Task{},
		&model.TaskHistory{},
		&model.Comment{},
		&model.Attachment{},
		&model.AuditLog{},
		&model.ProjectPermission{},
		&model.RolePermission{},
		&model.UserGroup{},
		&model.UserProjectRole{},
		&model.ProjectUser{},
		&model.TaskUser{},
	)
	if err != nil {
		panic("Failed to auto migrate database")
	}
	return db
}

func setupRoutes(app *fiber.App, db *gorm.DB) {
	router.Routes(app, db)
	app.Use(utils.NotFoundHandler)
}

func startServer(app *fiber.App, address string, errs chan<- error) {
	if err := app.Listen(address); err != nil {
		errs <- fmt.Errorf("error starting server: %w", err)
	}
}

func closeDatabase(db *gorm.DB) {
	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Error getting database instance: %v", errDB)
		return
	}

	if err := sqlDB.Close(); err != nil {
		utils.Log.Errorf("Error closing database connection: %v", err)
	} else {
		utils.Log.Info("Database connection closed successfully")
	}
}

func handleGracefulShutdown(ctx context.Context, app *fiber.App, serverErrors <-chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		utils.Log.Fatalf("Server error: %v", err)
	case <-quit:
		utils.Log.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			utils.Log.Fatalf("Error during server shutdown: %v", err)
		}
	case <-ctx.Done():
		utils.Log.Info("Server exiting due to context cancellation")
	}

	utils.Log.Info("Server exited")
}
