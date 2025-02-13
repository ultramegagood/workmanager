package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

// Хранилище WebSocket клиентов {userID: connection}
var clients = make(map[uuid.UUID]*websocket.Conn)
var lock sync.Mutex

// Структура WebSocket-сообщения
type WSMessage struct {
	TaskID      uuid.UUID `json:"task_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Comment     string    `json:"comment,omitempty"`
	UserID      uuid.UUID `json:"user_id"`
}

// Добавление клиента в WebSocket (с автоматическим удалением при закрытии)
func AddClient(userID uuid.UUID, conn *websocket.Conn) {
	lock.Lock()
	clients[userID] = conn
	lock.Unlock()

	// Удаление клиента при закрытии соединения
	defer RemoveClient(userID)

	// Чтение входящих сообщений (если нужно)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket closed for user %s: %v", userID, err)
			break
		}
	}
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

	// Создадим быстрый lookup map для O(1) поиска
	allowedMap := make(map[uuid.UUID]struct{})
	for _, userID := range allowedUsers {
		allowedMap[userID] = struct{}{}
	}

	// Пройдёмся по всем клиентам, а не по списку allowedUsers
	lock.Lock()
	defer lock.Unlock()
	for userID, conn := range clients {
		if _, allowed := allowedMap[userID]; allowed {
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("WebSocket send error:", err)
				RemoveClient(userID) // Если ошибка, удаляем клиента
			}
		}
	}
}
