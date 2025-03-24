package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/personal/task-management/internal/delivery/rest/dtos"
	_ "github.com/personal/task-management/internal/domain/task"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/apperrors"
	"github.com/personal/task-management/pkg/utils/jwt"
)

type TaskHandler struct {
	taskService usecase.TaskService
}

func NewTaskHandler(taskService usecase.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// godoc CreateTask
// @Summary Create Task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param createTaskInput body dtos.CreateTaskInput true "Create task input"
// @Success 201 {object} task.Task "Create task response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task dtos.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError(err.Error()))
		return
	}

	// get user id from context
	if userID, ok := r.Context().Value("user").(*jwt.UserClaims); ok {
		task.CreatorID = userID.UserID
	} else {
		apperrors.WriteError(w, apperrors.NewBadRequestError("User not found in context"))
		return
	}

	createdTask, err := h.taskService.CreateTask(r.Context(), task)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTask)
}

// godoc ListTasks
// @Summary List Tasks
// @Description List all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []task.Task "List tasks response"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	var input dtos.GetTasksWithFilterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError(err.Error()))
		return
	}

	tasks, err := h.taskService.GetTasksWithFilter(r.Context(), input)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// godoc GetEmployeeTasks
// @Summary Get Employee Tasks
// @Description Get tasks assigned to an employee
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Employee ID"
// @Success 200 {object} []task.Task "Get employee tasks response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 404 {object} apperrors.AppError "Not Found"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks/employee/{id} [get]
func (h *TaskHandler) GetEmployeeTasks(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	var requesterID uuid.UUID
	if userID, ok := r.Context().Value("user").(*jwt.UserClaims); ok {
		requesterID = userID.UserID
	} else {
		apperrors.WriteError(w, apperrors.NewBadRequestError("User not found in context"))
		return
	}

	employeeID := chi.URLParam(r, "id")
	employeeIDUUID, err := uuid.Parse(employeeID)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid employee ID"))
		return
	}

	input := dtos.GetEmployeeTasksInput{
		EmployeeID:  employeeIDUUID,
		RequesterID: requesterID,
	}

	tasks, err := h.taskService.GetEmployeeTasks(r.Context(), input)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// godoc GetSummaryByEmployee
// @Summary Get Summary By Employee
// @Description Get summary of tasks by employee
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []dtos.EmployeeTaskSummary "Get summary by employee response"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks/summary [get]
func (h *TaskHandler) GetSummaryByEmployee(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	var requesterID uuid.UUID
	if userID, ok := r.Context().Value("user").(*jwt.UserClaims); ok {
		requesterID = userID.UserID
	} else {
		apperrors.WriteError(w, apperrors.NewBadRequestError("User not found in context"))
		return
	}

	input := dtos.GetTaskSummaryByEmployeeInput{
		RequesterID: requesterID,
	}

	summary, err := h.taskService.GetTaskSummaryByEmployee(r.Context(), input)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// godoc GetTask
// @Summary Get Task
// @Description Get a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} task.Task "Get task response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 404 {object} apperrors.AppError "Not Found"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	var requesterID uuid.UUID
	if userID, ok := r.Context().Value("user").(*jwt.UserClaims); ok {
		requesterID = userID.UserID
	} else {
		apperrors.WriteError(w, apperrors.NewBadRequestError("User not found in context"))
		return
	}

	taskID := chi.URLParam(r, "id")
	taskIDUUID, err := uuid.Parse(taskID)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid task ID"))
		return
	}

	input := dtos.GetTaskInput{
		TaskID:      taskIDUUID,
		RequesterID: requesterID,
	}

	task, err := h.taskService.GetTask(r.Context(), input)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// godoc UpdateTask
// @Summary Update Task
// @Description Update a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param updateTaskInput body dtos.UpdateTaskStatusInput true "Update task input"
// @Success 200 {object} task.Task "Update task response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 404 {object} apperrors.AppError "Not Found"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input dtos.UpdateTaskStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError(err.Error()))
		return
	}

	// get user id from context
	if userID, ok := r.Context().Value("user").(*jwt.UserClaims); ok {
		input.UserID = userID.UserID
	} else {
		apperrors.WriteError(w, apperrors.NewBadRequestError("User not found in context"))
		return
	}
	task, err := h.taskService.UpdateTaskStatus(r.Context(), input)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// godoc DeleteTask
// @Summary Delete Task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} task.Task "Delete task response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 404 {object} apperrors.AppError "Not Found"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	var input dtos.DeleteTaskInput
	if userID, ok := r.Context().Value("user").(*jwt.UserClaims); ok {
		input = dtos.DeleteTaskInput{
			RequesterID: userID.UserID,
			TaskID:      uuid.MustParse(taskID),
		}
	} else {
		apperrors.WriteError(w, apperrors.NewBadRequestError("User not found in context"))
		return
	}

	err := h.taskService.DeleteTask(r.Context(), input)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}
