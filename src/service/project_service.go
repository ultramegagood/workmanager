package service

import (
	"app/src/model"
	"app/src/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateProject(c *fiber.Ctx, req *validation.CreateProject) (*model.Project, error)
	CreateTask(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error)
	CreateGroup(c *fiber.Ctx, req *validation.CreateGroup) (*model.Group, error)
	GetUsersWithAccess(taskID uuid.UUID) []uuid.UUID
	CreateUserGroup(c *fiber.Ctx, req *validation.CreateUserGroup) (*model.UserGroup, error)
	AddUserToGroup(c *fiber.Ctx, req *validation.AddUserToGroup) error
	AddGroupToProject(c *fiber.Ctx, req *validation.AddGroupToProject) error
	AddGroupToTask(c *fiber.Ctx, req *validation.AddGroupToTask) error
	GetUserGroups(c *fiber.Ctx) ([]model.UserGroup, error)
	GetUsersInGroup(c *fiber.Ctx, req *validation.GetUsersInGroup) ([]model.User, error)
}

type taskService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewTaskService(db *gorm.DB, validate *validator.Validate) TaskService {
	return &taskService{
		Log:      logrus.New(),
		DB:       db,
		Validate: validate,
	}
}
func (s *taskService) CreateProject(c *fiber.Ctx, req *validation.CreateProject) (*model.Project, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	project := &model.Project{
		Title: req.Title,
	}

	if err := s.DB.WithContext(c.Context()).Create(project).Error; err != nil {
		s.Log.Errorf("Failed to create project: %+v", err)
		return nil, err
	}

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
