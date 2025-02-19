package controller

import (
	"app/src/config"
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/utils"
	"app/src/validation"
	"strings"

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
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	project, err := tc.TaskService.CreateProject(c, &req, uuid.MustParse(userID))
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
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	var req validation.CreateTask
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	task, err := tc.TaskService.CreateTask(c, &req, uuid.MustParse(userID))
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

// CreateSection creates a new group.
// @Summary Create a new section
// @Description Create a new section with the provided details.
// @Tags Sections
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.CreateGroup true "Project section creation request"
// @Success 200 {object} response.SuccessWithData[model.Section]
// @Failure 400 {object} response.ErrorResponse
// @Router /projects/section [post]
func (tc *TaskController) CreateSection(c *fiber.Ctx) error {
	var req validation.CreateGroup
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	group, err := tc.TaskService.CreateProjectSection(c, &req)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.Section]{
		Code:    200,
		Status:  "success",
		Message: "Section added successfully",
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
// @Success 200 {object} response.SuccessWithData[[]model.User]
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks/{taskID}/users [get]
func (tc *TaskController) GetUsersWithAccess(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid task ID")
	}
	users := tc.TaskService.GetUsersWithAccess(taskID)
	return c.JSON(response.SuccessWithData[[]model.User]{
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
// @Description Добавление группы пользователей на задачу.
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

// Get tasks of user.
// @Summary Get tasks of user
// @Description Retrieve a list of tasks.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Success 200 {object} response.SuccessWithPaginate[model.Task]
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks [get]
func (tc *TaskController) GetTasks(c *fiber.Ctx) error {
	var tasks []model.Task
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	tasks, err = tc.TaskService.GetUserTasks(uuid.MustParse(userID))
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithPaginate[model.Task]{
		Code:    200,
		Status:  "success",
		Message: "Users in group retrieved successfully",
		Results: tasks,
	})
}

// GetUsersInGroup retrieves users in a specific group.
// @Summary Get users in a group
// @Description Retrieve a list of users in a specific group.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.GetUsersInGroup true "Get users in group request"
// @Success 200 {object} response.SuccessWithPaginate[model.User]
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks [put]
func (tc *TaskController) UpdateTaskTitleOrDescription(c *fiber.Ctx) error {

	var req validation.UpdateTaskTitleOrDescription
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err := tc.TaskService.UpdateTaskTitleOrDescription(c, req.TaskID, req.Title, req.Description)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.User]{
		Code:    200,
		Status:  "success",
		Message: "Task updated successfully",
	})
}

// GetProjects retrieves users in a specific group.
// @Summary Get projects by user
// @Description Retrieve a list of projects of user.
// @Tags Projects
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Success 200 {object} response.SuccessWithPaginate[model.Project]
// @Failure 400 {object} response.ErrorResponse
// @Router /projects [get]
func (tc *TaskController) GetUserProjects(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	projects, err := tc.TaskService.GetUserProjects(uuid.MustParse(userID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch projects")
	}

	return c.JSON(response.SuccessWithPaginate[model.Project]{
		Code:    200,
		Status:  "success",
		Message: "Projects fetched successfully",
		Results: projects, // Была ошибка в имени переменной
	})
}

// Comment task.
// @Summary Comment task
// @Tags Comments
// @Accept json
// @Produce json
// @Security  BearerAuth
// @Param request body validation.CreateComment true "Create comment"
// @Success 200 {object} response.SuccessWithData[model.Comment]
// @Failure 400 {object} response.ErrorResponse
// @Router /comments [post]
func (tc *TaskController) CommentTask(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	var req validation.CreateComment
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	comment, err := tc.TaskService.CreateComment(c, &req, uuid.MustParse(userID))
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[*model.Comment]{
		Code:    200,
		Status:  "success",
		Message: "Comment send successfully",
		Data:    comment,
	})

}

// ReassignTask reassigns a task to a new user.
// @Summary Reassign task to a new user
// @Description Change the assignee of a task and notify via WebSocket.
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Param request body validation.ReassignTaskValidation true "New assignee information"
// @Success 200 {object} response.Common
// @Failure 400 {object} response.ErrorResponse
// @Router /tasks/{taskID}/reassign [put]
func (tc *TaskController) ReassignTask(c *fiber.Ctx) error {
	var req validation.ReassignTaskValidation
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	if err := tc.TaskService.ReassignTask(c, req); err != nil {
		return err
	}

	return c.JSON(response.Common{
		Code:    200,
		Status:  "success",
		Message: "Task successfully reassigned",
	})
}

// Get task by ID.
// @Summary Get task by ID
// @Description Retrieve a task by its unique ID.
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Success 200 {object} response.SuccessWithData[model.Task]
// @Failure 404 {object} response.ErrorResponse
// @Router /tasks/{taskID} [get]
func (tc *TaskController) GetTaskByID(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid task ID")
	}
	task, err := tc.TaskService.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithData[model.Task]{
		Code:    200,
		Status:  "success",
		Message: "Task retrieved successfully",
		Data:    *task,
	})
}

// Delete task by ID.
// @Summary Delete task by ID
// @Description Delete a task by its unique ID.
// @Tags Tasks
// @Security BearerAuth
// @Param taskID path string true "Task ID"
// @Success 200 {object} response.Common
// @Failure 404 {object} response.ErrorResponse
// @Router /tasks/{taskID} [delete]
func (tc *TaskController) DeleteTask(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid task ID")
	}
	if err := tc.TaskService.DeleteTask(taskID); err != nil {
		return err
	}
	return c.JSON(response.Common{
		Code:    200,
		Status:  "success",
		Message: "Task deleted successfully",
	})
}

// Get sections by project ID.
// @Summary Get sections of a project
// @Description Retrieve all sections within a specific project.
// @Tags Sections
// @Produce json
// @Security BearerAuth
// @Param projectID path string true "Project ID"
// @Success 200 {object} response.SuccessWithPaginate[model.Section]
// @Failure 404 {object} response.ErrorResponse
// @Router /projects/{projectID}/sections [get]
func (tc *TaskController) GetSectionsByProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid project ID")
	}
	sections, err := tc.TaskService.GetSectionsByProject(projectID)
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithPaginate[model.Section]{
		Code:    200,
		Status:  "success",
		Message: "Sections retrieved successfully",
		Results: sections,
	})
}

// Get sections by user ID.
// @Summary Get sections of user
// @Description Retrieve all sections within a specific user.
// @Tags Sections
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessWithPaginate[model.UserSection]
// @Failure 404 {object} response.ErrorResponse
// @Router /sections [get]
func (tc *TaskController) GetSectionsByUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
	}

	userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
	if err != nil {
		return err
	}
	sections, err := tc.TaskService.GetSectionsByUser(uuid.MustParse(userID))
	if err != nil {
		return err
	}
	return c.JSON(response.SuccessWithPaginate[model.UserSection]{
		Code:    200,
		Status:  "success",
		Message: "Sections retrieved successfully",
		Results: sections,
	})
}

// Delete section by ID.
// @Summary Delete section by ID
// @Description Delete a section by its unique ID.
// @Tags Sections
// @Security BearerAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} response.Common
// @Failure 404 {object} response.ErrorResponse
// @Router /sections/{sectionID} [delete]
func (tc *TaskController) DeleteSection(c *fiber.Ctx) error {
	sectionID, err := uuid.Parse(c.Params("sectionID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid section ID")
	}
	if err := tc.TaskService.DeleteSection(sectionID); err != nil {
		return err
	}
	return c.JSON(response.Common{
		Code:    200,
		Status:  "success",
		Message: "Section deleted successfully",
	})
}
