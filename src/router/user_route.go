package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func UserRoutes(v1 fiber.Router, u service.UserService, t service.TokenService, task service.TaskService) {
	userController := controller.NewUserController(u, t)
	user := v1.Group("/users")

	// Обычные REST API-эндпоинты
	user.Get("/", m.Auth(u, "getUsers"), userController.GetUsers)
	user.Post("/", m.Auth(u, "manageUsers"), userController.CreateUser)
	user.Get("/:userId", m.Auth(u, "getUsers"), userController.GetUserByID)
	user.Patch("/:userId", m.Auth(u, "manageUsers"), userController.UpdateUser)
	user.Delete("/:userId", m.Auth(u, "manageUsers"), userController.DeleteUser)

	// WebSocket маршрут (выносится отдельно!)
	user.Get("/ws/:user_id", websocket.New(func(conn *websocket.Conn) {
		userID, err := uuid.Parse(conn.Params("user_id"))
		if err != nil {
			log.Println("Invalid user ID")
			conn.Close()
			return
		}

		// Добавляем клиента в WebSocket-менеджер
		service.AddClient(userID, conn)
		defer service.RemoveClient(userID)

		// Закрытие соединения обработчик
		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("WebSocket closed: user %s, code %d, message %s", userID, code, text)
			service.RemoveClient(userID)
			return nil
		})

		// Чтение сообщений от клиента
		for {
			var msg service.WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				log.Println("Error reading message:", err)
				break // Разрываем соединение
			}

			// Получаем пользователей с доступом к задаче
			allowedUsers := task.GetUsersWithAccess(msg.TaskID)

			// Отправляем сообщение только разрешённым пользователям
			service.BroadcastUpdate(msg, allowedUsers)
		}
	}))
}
