package controller

import (
	"app/src/config"
	"app/src/service"
	"app/src/utils"
	"log"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

// WebSocket обработчик
func WebSocketHandler(c *websocket.Conn) {
	// Получаем токен из query параметра (?token=xxx)
	token := c.Query("token")
	if token == "" {
		log.Println("Missing token in WebSocket connection")
		c.Close()
		return
	}

	// Проверяем токен и получаем userID
	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		log.Println("Invalid WebSocket token:", err)
		c.Close()
		return
	}

	log.Printf("User %s connected to WebSocket", userID)
	service.AddClient(uuid.MustParse(userID), c)
}
