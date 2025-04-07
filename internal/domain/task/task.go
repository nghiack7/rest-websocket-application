package task

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the current state of a task
type Status string

const (
	// StatusPending represents a task that has not been started
	StatusPending Status = "pending"
	// StatusInProgress represents a task that is currently being worked on
	StatusInProgress Status = "in_progress"
	// StatusCompleted represents a task that has been completed
	StatusCompleted Status = "completed"
	// StatusDeleted represents a task that has been deleted
	StatusDeleted Status = "deleted"
)

func (s Status) String() string {
	return string(s)
}

// Task represents a task in the system
type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	AssigneeID  uuid.UUID `json:"assignee_id"`
	CreatorID   uuid.UUID `json:"creator_id"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTask creates a new task with the given parameters
func NewTask(title, description string, dueDate time.Time, creatorID, assigneeID uuid.UUID) (*Task, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	if dueDate.Before(time.Now()) {
		return nil, ErrInvalidDueDate
	}

	now := time.Now()
	return &Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      StatusPending, // Default status for new tasks
		AssigneeID:  assigneeID,
		CreatorID:   creatorID,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// UpdateStatus updates the task status if the transition is valid
func (t *Task) UpdateStatus(newStatus Status) error {
	if !isValidStatusTransition(t.Status, newStatus) {
		return ErrInvalidStatusTransition
	}

	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(current, next Status) bool {
	// Define valid transitions
	switch current {
	case StatusPending:
		return next == StatusInProgress || next == StatusCompleted
	case StatusInProgress:
		return next == StatusCompleted
	case StatusCompleted:
		return false // Completed tasks cannot transition to other statuses
	default:
		return false
	}
}

// IsAssignedTo checks if the task is assigned to the given user
func (t *Task) IsAssignedTo(userID uuid.UUID) bool {
	return t.AssigneeID == userID
}

// IsCreatedBy checks if the task was created by the given user
func (t *Task) IsCreatedBy(userID uuid.UUID) bool {
	return t.CreatorID == userID
}

// IsPending checks if the task is in pending status
func (t *Task) IsPending() bool {
	return t.Status == StatusPending
}

// IsInProgress checks if the task is in progress
func (t *Task) IsInProgress() bool {
	return t.Status == StatusInProgress
}

// IsCompleted checks if the task is completed
func (t *Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}
