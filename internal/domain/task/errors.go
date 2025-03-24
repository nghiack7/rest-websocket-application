package task

import "errors"

// Task domain errors
var (
	ErrEmptyTitle              = errors.New("title cannot be empty")
	ErrInvalidDueDate          = errors.New("due date must be in the future")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrTaskNotFound            = errors.New("task not found")
	ErrUnauthorized            = errors.New("unauthorized to perform this action on the task")
)
