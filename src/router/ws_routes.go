package router

import (
	"app/src/controller"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupWebSocketRoutes(app *fiber.App) {
	// WebSocket маршрут
	app.Get("/ws", websocket.New(controller.WebSocketHandler))
}
