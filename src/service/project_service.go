package service

import (
	"app/src/model"
	"app/src/validation"
	"context"
	"encoding/json"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateProject(c *fiber.Ctx, req *validation.CreateProject, userID uuid.UUID) (*model.Project, error)
	CreateTask(c *fiber.Ctx, req *validation.CreateTask, userId uuid.UUID) (*model.Task, error)
	CreateProjectSection(c *fiber.Ctx, req *validation.CreateGroup) (*model.Section, error)
	GetUsersWithAccess(taskID uuid.UUID) []model.User
	CreateUserGroup(c *fiber.Ctx, req *validation.CreateUserGroup) (*model.UserGroup, error)
	AddUserToGroup(c *fiber.Ctx, req *validation.AddUserToGroup) error
	AddGroupToProject(c *fiber.Ctx, req *validation.AddGroupToProject) error
	AddGroupToTask(c *fiber.Ctx, req *validation.AddGroupToTask) error
	GetUserGroups(c *fiber.Ctx) ([]model.UserGroup, error)
	GetUsersInGroup(c *fiber.Ctx, req *validation.GetUsersInGroup) ([]model.User, error)
	HandleTaskUpdates(c *websocket.Conn)
	HandleProjectUpdates(c *websocket.Conn)
	GetUserTasks(userID uuid.UUID) ([]model.Task, error)
	HandleCommentUpdates(c *websocket.Conn)
	GetUserProjects(userID uuid.UUID) ([]model.Project, error)
	UpdateTaskTitleOrDescription(c *fiber.Ctx, taskID uuid.UUID, title, description string) error
	CreateComment(c *fiber.Ctx, req *validation.CreateComment, userID uuid.UUID) (*model.Comment, error)
	ReassignTask(c *fiber.Ctx, req validation.ReassignTaskValidation) error
	GetTaskByID(taskID uuid.UUID) (*model.Task, error)
	DeleteTask(taskID uuid.UUID) error
	GetSectionsByProject(projectID uuid.UUID) ([]model.Section, error)
	GetSectionsByUser(userID uuid.UUID) ([]model.UserSection, error)
	DeleteSection(sectionID uuid.UUID) error

}

func NewTaskService(db *gorm.DB, validate *validator.Validate, redisClient *redis.Client) TaskService {
	return &taskService{
		Log:      logrus.New(),
		DB:       db,
		Validate: validate,
		Redis:    redisClient,
	}
}

type taskService struct {
	Log       *logrus.Logger
	DB        *gorm.DB
	Validate  *validator.Validate
	Redis     *redis.Client
	WebSocket *websocket.Conn
}


const (
	debounceDuration      = 1 * time.Second
    taskUpdatesChannel    = "task_updates"
    commentUpdatesChannel = "comment_updates"
    projectUpdatesChannel = "project_updates"
)

// Общий тип для WebSocket сообщений
type WSMessage struct {
    Entity    string      `json:"entity"`   // "task", "comment", "project"
    Action    string      `json:"action"`   // "created", "updated", "deleted"
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}
func (s *taskService) CreateProject(c *fiber.Ctx, req *validation.CreateProject, userID uuid.UUID) (*model.Project, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	project := &model.Project{
		Title: req.Title,
	}

	tx := s.DB.WithContext(c.Context()).Begin()
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to create project: %+v", err)
		return nil, err
	}

	// Добавляем пользователя в проект
	projectUser := &model.ProjectUser{
		ProjectID: project.ID,
		UserID:    userID,
	}
	if err := tx.Create(projectUser).Error; err != nil {
		tx.Rollback()
		s.Log.Errorf("Failed to add user to project: %+v", err)
		return nil, err
	}

	tx.Commit()
	return project, nil
}
func (s *taskService) CreateTask(c *fiber.Ctx, req *validation.CreateTask, userID uuid.UUID) (*model.Task, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	var project model.Project
	if err := s.DB.First(&project, "id = ?", req.ProjectID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	// Назначаем создателя задачи, если не передан другой пользователь
	assignedTo := req.AssignedTo
	if assignedTo == nil {
		assignedTo = &userID
	}

	// Ищем секцию "Recently Assigned" пользователя
	var userSection model.UserSection
	if err := s.DB.
		Where("user_id = ? AND title = ?", *assignedTo, "Recently Assigned").
		First(&userSection).Error; err != nil {
		s.Log.Warnf("UserSection 'Recently Assigned' not found for user %s", assignedTo)
		return nil, fiber.NewError(fiber.StatusNotFound, "User section not found")
	}

	// Создаем таск
	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		AssignedTo:  assignedTo,
		SectionID:   userSection.ID, // Назначаем в первую секцию
	}

	if err := s.DB.WithContext(c.Context()).Create(task).Error; err != nil {
		s.Log.Errorf("Failed to create task: %+v", err)
		return nil, err
	}

	// Отправка WebSocket-обновления
	go s.publishUpdate(context.Background(), commentUpdatesChannel, WSMessage{
        Entity:    "task",
        Action:    "created",
        Data:      task,
        Timestamp: time.Now(),
    })
	return task, nil
}

// Функция переназначения таска
func (s *taskService) ReassignTask(c *fiber.Ctx, req validation.ReassignTaskValidation) error {
	var task model.Task
	if err := s.DB.First(&task, "id = ?", req.TaskID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Task not found")
	}

	// Ищем новую секцию "Recently Assigned" для нового исполнителя
	var userSection model.UserSection
	if err := s.DB.
		Where("user_id = ? AND title = ?", req.NewUserID, "Recently Assigned").
		First(&userSection).Error; err != nil {
		s.Log.Warnf("UserSection 'Recently Assigned' not found for user %s", req.NewUserID)
		return fiber.NewError(fiber.StatusNotFound, "User section not found")
	}

	// Обновляем исполнителя и секцию
	task.AssignedTo = &req.NewUserID
	task.SectionID = userSection.ID

	if err := s.DB.WithContext(c.Context()).Save(&task).Error; err != nil {
		s.Log.Errorf("Failed to reassign task: %+v", err)
		return err
	}

	// Отправка WebSocket-сообщения
	go s.publishUpdate(context.Background(), commentUpdatesChannel, WSMessage{
        Entity:    "task",
        Action:    "reassigned",
        Data:      task,
        Timestamp: time.Now(),
    })
	return nil
}
func (s *taskService) CreateProjectSection(c *fiber.Ctx, req *validation.CreateGroup) (*model.Section, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Проверяем, существует ли проект
	var project model.Project
	if err := s.DB.First(&project, "id = ?", req.ProjectID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	section := &model.Section{
		Title:     req.Title,
		ProjectID: req.ProjectID,
	}
	if err := s.DB.WithContext(c.Context()).Create(section).Error; err != nil {
		s.Log.Errorf("Failed to create section: %+v", err)
		return nil, err
	}
	return section, nil

}

func (s *taskService) GetUsersWithAccess(taskID uuid.UUID) []model.User {
	var users []model.User

	err := s.DB.Raw(`
        SELECT DISTINCT u.* 
        FROM users u
        INNER JOIN user_groups ug ON u.id = ug.user_id
        INNER JOIN tasks t ON ug.id = t.user_group
        WHERE t.id = ?
    `, taskID).Scan(&users).Error

	if err != nil {
		s.Log.Errorf("DB error in GetUsersWithAccess: %+v", err)
		return nil
	}

	return users
}

func (s *taskService) GetUserTasks(userID uuid.UUID) ([]model.Task, error) {
	var tasks []model.Task

	// Проверяем, существует ли пользователь
	var user model.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Получаем задачи, связанные с этим пользователем
	if err := s.DB.
		Joins("JOIN user_tasks ON user_tasks.task_id = tasks.id").
		Where("user_tasks.user_id = ?", userID).
		Find(&tasks).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	return tasks, nil
}

func (s *taskService) CreateUserGroup(c *fiber.Ctx, req *validation.CreateUserGroup) (*model.UserGroup, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	userGroup := &model.UserGroup{
		TeamTitle: req.TeamTitle,
		OwnerID:   req.UserID,
	}

	if err := s.DB.WithContext(c.Context()).Create(userGroup).Error; err != nil {
		s.Log.Errorf("Failed to create user group: %+v", err)
		return nil, err
	}

	return userGroup, nil
}

func (s *taskService) AddUserToGroup(c *fiber.Ctx, req *validation.AddUserToGroup) error {
	// Проверяем, существует ли группа
	s.Log.Infof("Adding user to group - UserID: %s, GroupID: %s", req.UserID, req.GroupID)
	var group model.UserGroup
	if err := s.DB.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Проверяем, существует ли пользователь
	var user model.User
	if err := s.DB.First(&user, "id = ?", req.UserID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Добавляем пользователя в группу
	if err := s.DB.Exec("INSERT INTO user_group_users (user_group_id, user_id) VALUES (?, ?) ON CONFLICT DO NOTHING", req.GroupID, req.UserID).Error; err != nil {
		s.Log.Errorf("Failed to add user to group: %+v", err)
		return err
	}

	return nil
}

func (s *taskService) AddGroupToProject(c *fiber.Ctx, req *validation.AddGroupToProject) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}
	// Проверяем, существует ли проект
	var project model.Project
	if err := s.DB.First(&project, "id = ?", req.ProjectID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	// Проверяем, существует ли группа
	var group model.UserGroup
	if err := s.DB.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Добавляем группу к проекту
	if err := s.DB.Model(&project).Association("UserGroups").Append(&group); err != nil {
		s.Log.Errorf("Failed to add group to project: %+v", err)
		return err
	}

	return nil
}
func (s *taskService) AddGroupToTask(c *fiber.Ctx, req *validation.AddGroupToTask) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}
	// Проверяем, существует ли задача
	var task model.Task
	if err := s.DB.First(&task, "id = ?", req.TaskID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Task not found")
	}

	// Проверяем, существует ли группа
	var group model.UserGroup
	if err := s.DB.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Добавляем группу к задаче
	if err := s.DB.Model(&task).Association("UserGroups").Append(&group); err != nil {

		s.Log.Errorf("Failed to add group to task: %+v", err)
		return err
	}

	return nil
}
func (s *taskService) GetUserGroups(c *fiber.Ctx) ([]model.UserGroup, error) {
	var userGroups []model.UserGroup

	if err := s.DB.WithContext(c.Context()).Find(&userGroups).Error; err != nil {
		s.Log.Errorf("Failed to get user groups: %+v", err)
		return nil, err
	}

	return userGroups, nil
}
func (s *taskService) GetUsersInGroup(c *fiber.Ctx, req *validation.GetUsersInGroup) ([]model.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	var group model.UserGroup
	if err := s.DB.Preload("Users").First(&group, "id = ?", req.GroupID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	return group.Users, nil
}
func (s *taskService) GetUserProjects(userID uuid.UUID) ([]model.Project, error) {

	var projects []model.Project
	err := s.DB.
		Joins("JOIN project_users ON project_users.project_id = projects.id").
		Where("project_users.user_id = ?", userID).
		Preload("Users"). // Загружаем пользователей в проекте
		Find(&projects).Error

	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Projects not found")
	}
	return projects, nil
}

func (s *taskService) UpdateTaskTitleOrDescription(c *fiber.Ctx, taskID uuid.UUID, title, description string) error {
	// Публикуем изменение в Redis
	err := s.Redis.Publish(c.Context(), taskUpdatesChannel, map[string]interface{}{
		"task_id":     taskID,
		"title":       title,
		"description": description,
		"updated_at":  time.Now(),
	}).Err()
	if err != nil {
		s.Log.Errorf("Failed to publish task update: %v", err)
		return err
	}

	// Запланировать сохранение в Postgres с дебаунсингом
	go s.debouncedSaveToPostgres(taskID, title, description)
	return nil
}

func (s *taskService) debouncedSaveToPostgres(taskID uuid.UUID, title, description string) {
	// Создаем уникальный ключ для дебаунсинга
	key := "task_debounce:" + taskID.String()

	// Устанавливаем время дебаунсинга
	s.Redis.SetNX(context.Background(), key, "1", debounceDuration)

	// Ждем указанное время перед сохранением
	time.Sleep(debounceDuration)

	// Проверяем, актуален ли еще наш запрос
	val, _ := s.Redis.GetDel(context.Background(), key).Result()
	if val == "" {
		return
	}

	// Сохраняем в Postgres
	err := s.DB.Model(&model.Task{}).
		Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"title":       title,
			"description": description,
			"updated_at":  time.Now(),
		}).Error

	if err != nil {
		s.Log.Errorf("Failed to save task updates: %v", err)
	}
}
// Общий обработчик WebSocket
func (s *taskService) HandleUpdates(c *websocket.Conn, channels ...string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    defer c.Close()

    // Подписываемся на все указанные каналы
    pubsub := s.Redis.Subscribe(ctx, channels...)
    defer pubsub.Close()

    ch := pubsub.Channel()

    // Горутина для чтения из WebSocket
    go func() {
        for {
            var msg WSMessage
            if err := c.ReadJSON(&msg); err != nil {
				s.Log.Errorf("Invalid WebSocket message format: %s", (err))
                continue
            }
            // Публикуем полученное сообщение в Redis
            if err := s.publishUpdate(ctx, msg.Entity+"_updates", msg); err != nil {
                s.Log.Errorf("Publish error: %v", err)
            }
        }
    }()

    // Основной цикл обработки сообщений из Redis
    for {
        select {
        case msg := <-ch:
            if err := c.WriteJSON(WSMessage{
                Timestamp: time.Now(),
                Data:      json.RawMessage(msg.Payload),
            }); err != nil {
                s.Log.Errorf("WebSocket write error:", err)
                continue
            }
        case <-ctx.Done():
            return
        }
    }
}
// Вспомогательный метод для публикации обновлений
func (s *taskService) publishUpdate(ctx context.Context, channel string, data interface{}) error {
    payload, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return s.Redis.Publish(ctx, channel, payload).Err()
}

// Обновленные обработчики для конкретных сущностей
func (s *taskService) HandleTaskUpdates(c *websocket.Conn) {
    s.HandleUpdates(c, taskUpdatesChannel)
}

func (s *taskService) HandleCommentUpdates(c *websocket.Conn) {
    s.HandleUpdates(c, commentUpdatesChannel)
}

func (s *taskService) HandleProjectUpdates(c *websocket.Conn) {
    s.HandleUpdates(c, projectUpdatesChannel)
}

func (s *taskService) CreateComment(c *fiber.Ctx, req *validation.CreateComment, userID uuid.UUID) (*model.Comment, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	comment := &model.Comment{
		Body:   req.Body,
		TaskID: req.TaskID,
		UserID: userID,
	}

	if err := s.DB.WithContext(c.Context()).Create(comment).Error; err != nil {
		s.Log.Errorf("Ошибка создания комментария: %+v", err)
		return nil, err
	}

	// Публикация комментария в Redis
	go s.publishUpdate(context.Background(), commentUpdatesChannel, WSMessage{
        Entity:    "comment",
        Action:    "created",
        Data:      comment,
        Timestamp: time.Now(),
    })

	return comment, nil
}

// Реализация в taskService:
func (s *taskService) GetTaskByID(taskID uuid.UUID) (*model.Task, error) {
	var task model.Task
	if err := s.DB.
		Preload("Comments").
		Preload("UserGroups").
		First(&task, "id = ?", taskID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Task not found")
	}
	return &task, nil
}

func (s *taskService) DeleteTask(taskID uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Удаляем связи задачи с группами
		if err := tx.Exec("DELETE FROM task_user_groups WHERE task_id = ?", taskID).Error; err != nil {
			return err
		}

		// Удаляем саму задачу
		if err := tx.Delete(&model.Task{}, "id = ?", taskID).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete task")
		}
		return nil
	})
}

func (s *taskService) GetSectionsByProject(projectID uuid.UUID) ([]model.Section, error) {
	var sections []model.Section
	if err := s.DB.
		Where("project_id = ?", projectID).
		Preload("Tasks").
		Find(&sections).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sections not found")
	}
	return sections, nil
}
func (s *taskService) GetSectionsByUser(userID uuid.UUID) ([]model.UserSection, error) {
	var sections []model.UserSection
	if err := s.DB.
		Where("user_id = ?", userID).
		Preload("Tasks").
		Find(&sections).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sections not found")
	}
	return sections, nil
}

func (s *taskService) DeleteSection(sectionID uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Переносим задачи в дефолтную секцию
		var defaultSection model.Section
		if err := tx.FirstOrCreate(&defaultSection,
			model.Section{Title: "Backlog", ProjectID: uuid.Nil}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Task{}).
			Where("section_id = ?", sectionID).
			Update("section_id", defaultSection.ID).Error; err != nil {
			return err
		}

		// Удаляем секцию
		if err := tx.Delete(&model.Section{}, "id = ?", sectionID).Error; err != nil {
			return err
		}

		return nil
	})
}
