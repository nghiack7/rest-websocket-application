package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/personal/task-management/internal/domain/task"
)

// TaskRepository defines the interface for task persistence operations
type TaskRepository interface {
	// Create stores a new task in the repository
	Create(ctx context.Context, task *task.Task) error

	// GetByID retrieves a task by ID
	GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error)

	// Update updates an existing task in the repository
	Update(ctx context.Context, task *task.Task) error

	// Delete removes a task from the repository
	Delete(ctx context.Context, id uuid.UUID) error

	// FindByAssignee retrieves tasks assigned to a specific user
	FindByAssignee(ctx context.Context, assigneeID uuid.UUID) ([]*task.Task, error)

	// FindByCreator retrieves tasks created by a specific user
	FindByCreator(ctx context.Context, creatorID uuid.UUID) ([]*task.Task, error)

	// FindByStatus retrieves tasks with a specific status
	FindByStatus(ctx context.Context, status task.Status) ([]*task.Task, error)

	// FindByDueDateRange retrieves tasks with due dates in a given range
	FindByDueDateRange(ctx context.Context, start, end time.Time) ([]*task.Task, error)

	// List retrieves all tasks with optional filtering and sorting
	List(ctx context.Context, filter TaskFilter) ([]*task.Task, error)
}

// TaskFilter defines filtering and sorting options for tasks
type TaskFilter struct {
	AssigneeID *uuid.UUID   `json:"assignee_id,omitempty"`
	Status     *task.Status `json:"status,omitempty"`
	SortBy     string       `json:"sort_by,omitempty"`    // Options: "due_date", "status", "created_at"
	SortOrder  string       `json:"sort_order,omitempty"` // Options: "asc", "desc"
	Offset     int          `json:"offset,omitempty"`
	Limit      int          `json:"limit,omitempty"`
}
