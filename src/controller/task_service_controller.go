package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TaskController struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

// CreateProject creates a new project.
// @Summary Create a new project
// @Description Create a new project with the provided details.
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body validation.CreateProject true "Project creation request"
// @Success 200 {object} response.SuccessWithData[model.Project]
// @Failure 400 {object} response.ErrorResponse
// @Router /projects [post]
func (tc *TaskController) CreateProject(c *fiber.Ctx) error {
	var req validation.CreateProject
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	project, err := tc.TaskService.CreateProject(c, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.Project]{
		Code:    200,
		Status:  "success",
		Message: "Project created successfully",
		Data:    *project,
	})
}

// CreateTask creates a new task.
// @Summary Create a new task
// @Description Create a new task with the provided details.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.CreateTask true "Task creation request"
// @Success 200 {object} response.SuccessWithData[model.Task]
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks [post]
func (tc *TaskController) CreateTask(c *fiber.Ctx) error {
	var req validation.CreateTask
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	task, err := tc.TaskService.CreateTask(c, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.Task]{
		Code:    200,
		Status:  "success",
		Message: "Task created successfully",
		Data:    *task,
	})
}

// CreateGroup creates a new group.
// @Summary Create a new group
// @Description Create a new group with the provided details.
// @Tags Groups
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.CreateGroup true "Group creation request"
// @Success 200 {object} response.SuccessWithData[model.Group]
// @Failure 400 {object} response.ErrorResponse
// @Router /groups [post]
func (tc *TaskController) CreateGroup(c *fiber.Ctx) error {
	var req validation.CreateGroup
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	group, err := tc.TaskService.CreateGroup(c, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.Group]{
		Code:    200,
		Status:  "success",
		Message: "Group created successfully",
		Data:    *group,
	})
}

// GetUsersWithAccess retrieves users with access to a task.
// @Summary Get users with access to a task
// @Description Retrieve a list of users who have access to a specific task.
// @Tags Tasks
// @Produce json
// @Security  BearerAuth
// @Param taskID path string true "Task ID"
// @Success 200 {object} response.SuccessWithData[[]uuid.UUID]
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks/{taskID}/users [get]
func (tc *TaskController) GetUsersWithAccess(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid task ID")
	}
	users := tc.TaskService.GetUsersWithAccess(taskID)
	return c.JSON(response.SuccessWithData[[]uuid.UUID]{
		Code:    200,
		Status:  "success",
		Message: "Users with access retrieved successfully",
		Data:    users,
	})
}

// CreateUserGroup creates a new user group.
// @Summary Create a new user group
// @Description Create a new user group with the provided details.
// @Tags UserGroups
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.CreateUserGroup true "User group creation request"
// @Success 200 {object} response.SuccessWithData[model.UserGroup]
// @Failure 400 {object} response.ErrorResponse
// @Router /user-groups [post]
func (tc *TaskController) CreateUserGroup(c *fiber.Ctx) error {
	var req validation.CreateUserGroup
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	group, err := tc.TaskService.CreateUserGroup(c, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.UserGroup]{
		Code:    200,
		Status:  "success",
		Message: "User group created successfully",
		Data:    *group,
	})
}

// AddUserToGroup adds a user to a group.
// @Summary Add a user to a group
// @Description Add a user to an existing group.
// @Tags UserGroups
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.AddUserToGroup true "Add user to group request"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.ErrorResponse
// @Router /user-groups/add-user [post]
func (tc *TaskController) AddUserToGroup(c *fiber.Ctx) error {
	var req validation.AddUserToGroup
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := tc.TaskService.AddUserToGroup(c, &req); err != nil {
		return err
	}
	return c.JSON(response.Common{
		Code:    200,
		Status:  "success",
		Message: "User added to group successfully",
	})
}

// AddGroupToProject adds a group to a project.
// @Summary Add a group to a project
// @Description Add a group to an existing project.
// @Tags Projects
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.AddGroupToProject true "Add group to project request"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.ErrorResponse
// @Router /projects/add-group [post]
func (tc *TaskController) AddGroupToProject(c *fiber.Ctx) error {
	var req validation.AddGroupToProject
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := tc.TaskService.AddGroupToProject(c, &req); err != nil {
		return err
	}
	return c.JSON(response.Common{
		Code:    200,
		Status:  "success",
		Message: "Group added to project successfully",
	})
}

// AddGroupToTask adds a group to a task.
// @Summary Add a group to a task
// @Description Add a group to an existing task.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.AddGroupToTask true "Add group to task request"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks/add-group [post]
func (tc *TaskController) AddGroupToTask(c *fiber.Ctx) error {
	var req validation.AddGroupToTask
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := tc.TaskService.AddGroupToTask(c, &req); err != nil {
		return err
	}
	return c.JSON(response.Common{
		Code:    200,
		Status:  "success",
		Message: "Group added to task successfully",
	})
}

// GetUserGroups retrieves all user groups.
// @Summary Get all user groups
// @Description Retrieve a list of all user groups.
// @Tags UserGroups
// @Produce json
// @Security  BearerAuth
// @Success 200 {object} response.SuccessWithPaginate[model.UserGroup]
// @Failure 400 {object} response.ErrorResponse
// @Router /user-groups [get]
func (tc *TaskController) GetUserGroups(c *fiber.Ctx) error {
	groups, err := tc.TaskService.GetUserGroups(c)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithPaginate[model.UserGroup]{
		Code:    200,
		Status:  "success",
		Message: "User groups retrieved successfully",
		Results: groups,
	})
}

// GetUsersInGroup retrieves users in a specific group.
// @Summary Get users in a group
// @Description Retrieve a list of users in a specific group.
// @Tags UserGroups
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.GetUsersInGroup true "Get users in group request"
// @Success 200 {object} response.SuccessWithPaginate[model.User]
// @Failure 400 {object} response.ErrorResponse
// @Router /user-groups/users [post]
func (tc *TaskController) GetUsersInGroup(c *fiber.Ctx) error {
	var req validation.GetUsersInGroup
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	users, err := tc.TaskService.GetUsersInGroup(c, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithPaginate[model.User]{
		Code:    200,
		Status:  "success",
		Message: "Users in group retrieved successfully",
		Results: users,
	})
}
