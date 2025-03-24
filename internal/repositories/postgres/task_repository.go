package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/personal/task-management/internal/domain/task"
	repository "github.com/personal/task-management/internal/repositories"
	"github.com/personal/task-management/pkg/cache"
	"gorm.io/gorm"
)

type PostgresTaskRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

func NewPostgresTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *task.Task) error {
	return r.db.Create(task).Error
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	var t task.Task
	if err := r.db.First(&t, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *task.Task) error {
	return r.db.Save(task).Error
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(&task.Task{}, "id = ?", id).Error
}

func (r *PostgresTaskRepository) FindByAssignee(ctx context.Context, assigneeID uuid.UUID) ([]*task.Task, error) {
	var tasks []*task.Task
	if err := r.db.Where("assignee_id = ?", assigneeID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) FindByCreator(ctx context.Context, creatorID uuid.UUID) ([]*task.Task, error) {
	var tasks []*task.Task
	if err := r.db.Where("creator_id = ?", creatorID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) FindByStatus(ctx context.Context, status task.Status) ([]*task.Task, error) {
	var tasks []*task.Task
	if err := r.db.Where("status = ?", status).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) FindByDueDateRange(ctx context.Context, start, end time.Time) ([]*task.Task, error) {
	var tasks []*task.Task
	if err := r.db.Where("due_date BETWEEN ? AND ?", start, end).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
func (r *PostgresTaskRepository) List(ctx context.Context, filter repository.TaskFilter) ([]*task.Task, error) {
	query := r.db.Model(&task.Task{})

	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", filter.AssigneeID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", filter.Status)
	}

	// Default sorting if not specified
	if filter.SortBy == "" {
		filter.SortBy = "created_at" // Default sort by creation date
		filter.SortOrder = "desc"    // Default newest first
	}

	// Handle specific sort fields
	switch filter.SortBy {
	case "created_at", "due_date":
		// These are valid date fields for sorting
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder))
	case "status":
		// Special handling for status sorting
		if filter.SortOrder == "asc" {
			// Order by status: pending, in_progress, completed
			query = query.Order("CASE status " +
				"WHEN 'pending' THEN 1 " +
				"WHEN 'in_progress' THEN 2 " +
				"WHEN 'completed' THEN 3 END")
		} else {
			// Order by status: completed, in_progress, pending
			query = query.Order("CASE status " +
				"WHEN 'completed' THEN 1 " +
				"WHEN 'in_progress' THEN 2 " +
				"WHEN 'pending' THEN 3 END")
		}
	default:
		// For any other fields, use standard ordering
		query = query.Order(fmt.Sprintf("%s %s", filter.SortBy, filter.SortOrder))
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	tasks := []*task.Task{}
	err := query.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
