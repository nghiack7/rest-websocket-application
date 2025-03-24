package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/personal/task-management/internal/domain/task"
)

type CreateTaskInput struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" validate:"required,gt=now"`
	AssigneeID  uuid.UUID `json:"assignee_id" validate:"required"`
	CreatorID   uuid.UUID `json:"creator_id" validate:"required"`
}

type UpdateTaskStatusInput struct {
	TaskID    uuid.UUID   `json:"task_id" validate:"required"`
	UserID    uuid.UUID   `json:"user_id" validate:"required"`
	NewStatus task.Status `json:"new_status" validate:"required,oneof=pending in_progress completed"`
}

type GetEmployeeTasksInput struct {
	EmployeeID  uuid.UUID `json:"employee_id" validate:"required"`
	RequesterID uuid.UUID `json:"requester_id" validate:"required"`
}

type GetTaskInput struct {
	TaskID      uuid.UUID `json:"task_id" validate:"required"`
	RequesterID uuid.UUID `json:"requester_id" validate:"required"`
}

type DeleteTaskInput struct {
	TaskID      uuid.UUID `json:"task_id" validate:"required"`
	RequesterID uuid.UUID `json:"requester_id" validate:"required"`
}

type GetTasksWithFilterInput struct {
	UserID uuid.UUID  `json:"user_id" validate:"required"`
	Filter TaskFilter `json:"filter" validate:"required"`
}

type TaskFilter struct {
	SortBy     string      `json:"sort_by"`
	Status     task.Status `json:"status"`
	DueDate    time.Time   `json:"due_date"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	SortOrder  string      `json:"sort_order"`
	AssigneeID uuid.UUID   `json:"assignee_id"`
}

type GetTaskSummaryByEmployeeInput struct {
	RequesterID uuid.UUID `json:"requester_id" validate:"required"`
}

type EmployeeTaskSummary struct {
	EmployeeID      uuid.UUID `json:"employee_id"`
	EmployeeName    string    `json:"employee_name"`
	TotalTasks      int       `json:"total_tasks"`
	CompletedTasks  int       `json:"completed_tasks"`
	PendingTasks    int       `json:"pending_tasks"`
	InProgressTasks int       `json:"in_progress_tasks"`
}
