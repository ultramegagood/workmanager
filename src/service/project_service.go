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
	CreateTask(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error)
	CreateGroup(c *fiber.Ctx, req *validation.CreateGroup) (*model.Group, error)
	GetUsersWithAccess(taskID uuid.UUID) []uuid.UUID
	CreateUserGroup(c *fiber.Ctx, req *validation.CreateUserGroup) (*model.UserGroup, error)
	AddUserToGroup(c *fiber.Ctx, req *validation.AddUserToGroup) error
	AddGroupToProject(c *fiber.Ctx, req *validation.AddGroupToProject) error
	AddGroupToTask(c *fiber.Ctx, req *validation.AddGroupToTask) error
	GetUserGroups(c *fiber.Ctx) ([]model.UserGroup, error)
	GetUsersInGroup(c *fiber.Ctx, req *validation.GetUsersInGroup) ([]model.User, error)
	HandleTaskUpdates(c *websocket.Conn)
	GetUserProjects(userID uuid.UUID) ([]model.Project, error)
	UpdateTaskTitleOrDescription(c *fiber.Ctx, taskID uuid.UUID, title, description string) error
}
type taskService struct {
	Log       *logrus.Logger
	DB        *gorm.DB
	Validate  *validator.Validate
	Redis     *redis.Client
	WebSocket *websocket.Conn
}

const (
	debounceDuration   = 1 * time.Second
	taskUpdatesChannel = "task_updates"
)

func NewTaskService(db *gorm.DB, validate *validator.Validate) TaskService {
	return &taskService{
		Log:      logrus.New(),
		DB:       db,
		Validate: validate,
	}
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


func (s *taskService) CreateTask(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Проверяем, существует ли проект
	var project model.Project
	if err := s.DB.First(&project, "id = ?", req.ProjectID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		ProjectID:   req.ProjectID,
	}

	if err := s.DB.WithContext(c.Context()).Create(task).Error; err != nil {
		s.Log.Errorf("Failed to create task: %+v", err)
		return nil, err
	}

	return task, nil
}

func (s *taskService) CreateGroup(c *fiber.Ctx, req *validation.CreateGroup) (*model.Group, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Проверяем, существует ли проект
	var project model.Project
	if err := s.DB.First(&project, "id = ?", req.ProjectID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Project not found")
	}

	group := &model.Group{
		Title:     req.Title,
		ProjectID: req.ProjectID,
	}

	if err := s.DB.WithContext(c.Context()).Create(group).Error; err != nil {
		s.Log.Errorf("Failed to create group: %+v", err)
		return nil, err
	}

	return group, nil
}
func (s *taskService) GetUsersWithAccess(taskID uuid.UUID) []uuid.UUID {
	var userIDs []uuid.UUID

	err := s.DB.Raw(`
        SELECT DISTINCT ug.user_id 
        FROM user_groups ug
        INNER JOIN tasks t ON ug.id = t.user_group
        WHERE t.id = ?
    `, taskID).Scan(&userIDs).Error

	if err != nil {
		s.Log.Errorf("DB error in GetUsersWithAccess: %+v", err)
	}

	return userIDs
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

	if(err !=nil){
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
func (s *taskService) UpdateTaskHandler(c *fiber.Ctx) error {
	var req struct {
		TaskID      uuid.UUID `json:"task_id" validate:"required"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	return s.UpdateTaskTitleOrDescription(c, req.TaskID, req.Title, req.Description)
}


func (s *taskService) HandleTaskUpdates(c *websocket.Conn) {
    ctx := context.Background()
    pubsub := s.Redis.Subscribe(ctx, taskUpdatesChannel)
    defer pubsub.Close()

    // Создаем канал для отслеживания закрытия соединения
    done := make(chan struct{})
    defer close(done)

    // Запускаем горутину для чтения сообщений
    go func() {
        for {
            select {
            case <-done:
                return
            default:
                _, _, err := c.ReadMessage()
                if err != nil {
                    // При ошибке чтения закрываем соединение
                    c.Close()
                    return
                }
            }
        }
    }()

    ch := pubsub.Channel()
    for {
        select {
        case msg := <-ch:
            err := c.WriteJSON(decodeMessage(msg.Payload))
            if err != nil {
                s.Log.Errorf("WebSocket write error: %v", err)
                return
            }
        case <-done:
            return
        }
    }
}

func decodeMessage(payload string) map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return nil
	}
	return data
}
