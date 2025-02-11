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
	GetUsersWithAccess(taskID uuid.UUID) ([]uuid.UUID)
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
