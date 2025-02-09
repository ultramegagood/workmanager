package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// Клиенты WebSocket {userID: connection}
var clients = make(map[uuid.UUID]*websocket.Conn)
var lock = sync.Mutex{}

// WebSocket сообщение
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
	clients[userID] = conn
	lock.Unlock()
}

// Удаление клиента
func RemoveClient(userID uuid.UUID) {
	lock.Lock()
	delete(clients, userID)
	lock.Unlock()
}

// Отправка обновлений всем пользователям, имеющим доступ
func BroadcastUpdate(msg WSMessage, allowedUsers []uuid.UUID) {
	data, _ := json.Marshal(msg)

	lock.Lock()
	for _, userID := range allowedUsers {
		if conn, exists := clients[userID]; exists {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("WebSocket send error:", err)
			}
		}
	}
	lock.Unlock()
}
