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

	user.Get("/", m.Auth(u, "getUsers"), userController.GetUsers)
	user.Post("/", m.Auth(u, "manageUsers"), userController.CreateUser)
	user.Get("/:userId", m.Auth(u, "getUsers"), userController.GetUserByID)
	user.Patch("/:userId", m.Auth(u, "manageUsers"), userController.UpdateUser)
	user.Delete("/:userId", m.Auth(u, "manageUsers"), userController.DeleteUser)
	user.Get("/ws/:user_id", websocket.New(func(c *websocket.Conn) {
		userID, err := uuid.Parse(c.Params("user_id"))
		if err != nil {
			log.Println("Invalid user ID")
			c.Close()
			return
		}
		service.AddClient(userID, c)
		defer service.RemoveClient(userID)
		for {
			var msg service.WSMessage
			if err := c.ReadJSON(&msg); err != nil {
				log.Println("Error reading message:", err)
				break
			}
			// Получаем список пользователей, имеющих доступ к задаче
			allowedUsers := task.GetUsersWithAccess(msg.TaskID)
			// Отправляем обновления только им
			service.BroadcastUpdate(msg, allowedUsers)
		}
	}))
}
