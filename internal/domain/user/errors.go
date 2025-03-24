package user

import "errors"

// User domain errors
var (
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrInvalidRole   = errors.New("invalid role")
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
)
