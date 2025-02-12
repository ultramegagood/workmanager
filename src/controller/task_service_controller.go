package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"

)

type TaskController struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

// @Tags         Tasks
// @Summary      Create a task
// @Description  Creates a new task in a specified project.
// @Security BearerAuth
// @Produce      json
// @Param        request  body  validation.CreateTask  true  "Request body"
// @Router       /tasks [post]
// @Success      201  {object}  example.CreateTaskResponse
// @Failure      400  {object}  example.BadRequest  "Invalid request"
// @Failure      404  {object}  example.NotFound  "Project not found"
func (t *TaskController) CreateTask(c *fiber.Ctx) error {
	req := new(validation.CreateTask)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	task, err := t.TaskService.CreateTask(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithTask{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Task created successfully",
			Task:    *task,
		})
}

// @Tags         Groups
// @Summary      Create a group
// @Description  Creates a new group within a project.
// @Security BearerAuth
// @Produce      json
// @Param        request  body  validation.CreateGroup  true  "Request body"
// @Router       /groups [post]
// @Success      201  {object}  example.CreateGroupResponse
// @Failure      400  {object}  example.BadRequest  "Invalid request"
// @Failure      404  {object}  example.NotFound  "Project not found"
func (t *TaskController) CreateGroup(c *fiber.Ctx) error {
	req := new(validation.CreateGroup)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	group, err := t.TaskService.CreateGroup(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithGroup{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Group created successfully",
			Group:   *group,
		})
}

// @Tags         Projects
// @Summary      Create a project
// @Description  Creates a new project.
// @Security BearerAuth
// @Produce      json
// @Param        request  body  validation.CreateProject  true  "Request body"
// @Router       /projects [post]
// @Success      201  {object}  example.CreateProjectResponse
// @Failure      400  {object}  example.BadRequest  "Invalid request"
func (t *TaskController) CreateProject(c *fiber.Ctx) error {
	req := new(validation.CreateProject)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	project, err := t.TaskService.CreateProject(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithProject{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Project created successfully",
			Project: *project,
		})
}
