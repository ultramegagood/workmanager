package service

import (
	"app/src/model"
	"app/src/validation"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskService interface {
	CreateProject(c *fiber.Ctx, req *validation.CreateProject) (*model.Project, error)
	CreateTask(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error)
	CreateTaskSupport(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error)
	CreateGroup(c *fiber.Ctx, req *validation.CreateGroup) (*model.Group, error)
	GetUsersWithAccess(taskID uuid.UUID) []uuid.UUID
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
		Title:    req.Title,
		GroupIDs: req.GroupIDs,
	}

	result := s.DB.WithContext(c.Context()).Create(project)
	if result.Error != nil {
		s.Log.Errorf("Failed to create project: %+v", result.Error)
		return nil, result.Error
	}

	return project, nil
}

func (s *taskService) CreateTask(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	task := &model.Task{}

	result := s.DB.WithContext(c.Context()).Create(task)
	if result.Error != nil {
		s.Log.Errorf("Failed to create task: %+v", result.Error)
		return nil, result.Error
	}

	return task, nil
}

func (s *taskService) CreateTaskSupport(c *fiber.Ctx, req *validation.CreateTask) (*model.Task, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		ProjectID:   uuid.MustParse(req.ProjectID),
	}

	result := s.DB.WithContext(c.Context()).Create(task)
	if result.Error != nil {
		s.Log.Errorf("Failed to create task: %+v", result.Error)
		return nil, result.Error
	}

	return task, nil
}

func (s *taskService) CreateGroup(c *fiber.Ctx, req *validation.CreateGroup) (*model.Group, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	group := &model.Group{
		Title: req.Title,
	}

	result := s.DB.WithContext(c.Context()).Create(group)
	if result.Error != nil {
		s.Log.Errorf("Failed to create group: %+v", result.Error)
		return nil, result.Error
	}

	return group, nil
}

func (s *taskService) GetUsersWithAccess(taskID uuid.UUID) []uuid.UUID {
	var userIDs []uuid.UUID

	// Запрос к БД: найти `user_group`, связанную с `taskID`
	err := s.DB.Raw(`
		SELECT user_id FROM user_groups 
		WHERE id IN (SELECT user_group FROM tasks WHERE id = ?)
	`, taskID).Scan(&userIDs).Error

	if err != nil {
		log.Println("DB error:", err)
	}

	return userIDs
}
