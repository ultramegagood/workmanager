package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// Хранилище WebSocket клиентов {userID: connection}
var clients = make(map[uuid.UUID]*websocket.Conn)
var lock sync.Mutex

// Структура сообщения WebSocket
type WSMessage struct {
	TaskID      uuid.UUID `json:"task_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Comment     string    `json:"comment,omitempty"`
	UserID      uuid.UUID `json:"user_id"`
}

// Добавление клиента в WebSocket
func AddClient(userID uuid.UUID, conn *websocket.Conn) {
	lock.Lock()
	defer lock.Unlock()
	clients[userID] = conn
}

// Удаление клиента
func RemoveClient(userID uuid.UUID) {
	lock.Lock()
	defer lock.Unlock()
	if conn, exists := clients[userID]; exists {
		conn.Close()
		delete(clients, userID)
	}
}

// Отправка обновлений только пользователям с доступом
func BroadcastUpdate(msg WSMessage, allowedUsers []uuid.UUID) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling WebSocket message:", err)
		return
	}

	lock.Lock()
	defer lock.Unlock()
	for _, userID := range allowedUsers {
		if conn, exists := clients[userID]; exists {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("WebSocket send error:", err)
				RemoveClient(userID) // Если ошибка, удаляем клиента
			}
		}
	}
}
