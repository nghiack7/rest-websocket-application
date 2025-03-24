package usecase

import (
	"context"

	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/domain/task"
	repository "github.com/personal/task-management/internal/repositories"
	"github.com/personal/task-management/pkg/utils/validate"
)

type TaskService interface {
	CreateTask(ctx context.Context, input dtos.CreateTaskInput) (*task.Task, error)
	UpdateTaskStatus(ctx context.Context, input dtos.UpdateTaskStatusInput) (*task.Task, error)
	GetTask(ctx context.Context, input dtos.GetTaskInput) (*task.Task, error)
	GetEmployeeTasks(ctx context.Context, input dtos.GetEmployeeTasksInput) ([]*task.Task, error)
	GetTasksWithFilter(ctx context.Context, input dtos.GetTasksWithFilterInput) ([]*task.Task, error)
	GetTaskSummaryByEmployee(ctx context.Context, input dtos.GetTaskSummaryByEmployeeInput) ([]dtos.EmployeeTaskSummary, error)
	DeleteTask(ctx context.Context, input dtos.DeleteTaskInput) error
}

// TaskService handles task-related operations and business logic
type taskService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
}

// NewTaskService creates a new instance of TaskService
func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// CreateTask creates a new task
func (s *taskService) CreateTask(ctx context.Context, input dtos.CreateTaskInput) (*task.Task, error) {
	// validate input
	err := validate.Struct(input)
	if err != nil {
		return nil, err
	}

	// Verify creator exists and has employer role
	creator, err := s.userRepo.GetByID(ctx, input.CreatorID)
	if err != nil {
		return nil, err
	}

	if !creator.CanCreateTasks() {
		return nil, task.ErrUnauthorized
	}

	// Verify assignee exists
	assignee, err := s.userRepo.GetByID(ctx, input.AssigneeID)
	if err != nil {
		return nil, err
	}

	if !assignee.IsEmployee() {
		return nil, task.ErrUnauthorized // Can only assign tasks to employees
	}

	// Create task
	newTask, err := task.NewTask(
		input.Title,
		input.Description,
		input.DueDate,
		input.CreatorID,
		input.AssigneeID,
	)
	if err != nil {
		return nil, err
	}

	// Save task
	if err := s.taskRepo.Create(ctx, newTask); err != nil {
		return nil, err
	}

	return newTask, nil
}

// UpdateTaskStatus updates the status of a task
func (s *taskService) UpdateTaskStatus(ctx context.Context, input dtos.UpdateTaskStatusInput) (*task.Task, error) {
	// Get task
	t, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if err != nil {
		return nil, err
	}

	// Get user
	u, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// Check authorization
	if !u.CanUpdateTaskStatus() {
		return nil, task.ErrUnauthorized
	}

	// Employees can only update tasks assigned to them
	if u.IsEmployee() && !t.IsAssignedTo(input.UserID) {
		return nil, task.ErrUnauthorized
	}

	// Update status
	if err := t.UpdateStatus(input.NewStatus); err != nil {
		return nil, err
	}

	// Save task
	if err := s.taskRepo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

// GetEmployeeTasks retrieves tasks assigned to an employee
func (s *taskService) GetEmployeeTasks(ctx context.Context, input dtos.GetEmployeeTasksInput) ([]*task.Task, error) {
	// Get requester
	requester, err := s.userRepo.GetByID(ctx, input.RequesterID)
	if err != nil {
		return nil, err
	}

	// Check authorization
	if requester.IsEmployee() && input.EmployeeID != input.RequesterID {
		return nil, task.ErrUnauthorized // Employees can only view their own tasks
	}
	// Get tasks
	return s.taskRepo.FindByAssignee(ctx, input.EmployeeID)
}

// GetTask retrieves a task by ID
func (s *taskService) GetTask(ctx context.Context, input dtos.GetTaskInput) (*task.Task, error) {
	// Get requester
	requester, err := s.userRepo.GetByID(ctx, input.RequesterID)
	if err != nil {
		return nil, err
	}

	// Get task
	t, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if err != nil {
		return nil, err
	}

	// Check authorization
	if requester.IsEmployee() && t.AssigneeID != input.RequesterID {
		return nil, task.ErrUnauthorized // Employees can only view their own tasks
	}

	return t, nil
}

// GetTasksWithFilter retrieves tasks with filtering and sorting
func (s *taskService) GetTasksWithFilter(ctx context.Context, input dtos.GetTasksWithFilterInput) ([]*task.Task, error) {
	// Get user
	u, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	// Check authorization for viewing all tasks
	if !u.CanViewAllTasks() {
		// Employee can only see their own tasks
		if u.IsEmployee() {
			input.Filter.AssigneeID = input.UserID
		}
	}
	filter := repository.TaskFilter{
		AssigneeID: &input.Filter.AssigneeID,
		Status:     &input.Filter.Status,
		Limit:      input.Filter.Limit,
		Offset:     input.Filter.Offset,
		SortBy:     input.Filter.SortBy,
		SortOrder:  input.Filter.SortOrder,
	}

	// Get tasks with filter
	return s.taskRepo.List(ctx, filter)
}

// GetTaskSummaryByEmployee retrieves a summary of tasks for all employees
func (s *taskService) GetTaskSummaryByEmployee(ctx context.Context, input dtos.GetTaskSummaryByEmployeeInput) ([]dtos.EmployeeTaskSummary, error) {
	// Get requester
	requester, err := s.userRepo.GetByID(ctx, input.RequesterID)
	if err != nil {
		return nil, err
	}

	// Only employers can see task summaries for all employees
	if !requester.IsEmployer() {
		return nil, task.ErrUnauthorized
	}

	// Get all employees
	users, err := s.userRepo.List(ctx, 0, 1000) // Pagination would be better in a real system
	if err != nil {
		return nil, err
	}

	var summaries []dtos.EmployeeTaskSummary
	for _, u := range users {
		if !u.IsEmployee() {
			continue // Only calculate summaries for employees
		}

		// Get tasks for this employee
		employeeTasks, err := s.taskRepo.FindByAssignee(ctx, u.ID)
		if err != nil {
			return nil, err
		}

		// Calculate summary
		summary := dtos.EmployeeTaskSummary{
			EmployeeID:   u.ID,
			EmployeeName: u.Name,
			TotalTasks:   len(employeeTasks),
		}

		// Count tasks by status
		for _, t := range employeeTasks {
			switch t.Status {
			case task.StatusPending:
				summary.PendingTasks++
			case task.StatusInProgress:
				summary.InProgressTasks++
			case task.StatusCompleted:
				summary.CompletedTasks++
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (s *taskService) DeleteTask(ctx context.Context, input dtos.DeleteTaskInput) error {
	// Get user
	u, err := s.userRepo.GetByID(ctx, input.RequesterID)
	if err != nil {
		return err
	}

	// Check authorization
	if !u.IsEmployer() {
		return task.ErrUnauthorized
	}

	// Delete task
	return s.taskRepo.Delete(ctx, input.TaskID)
}
