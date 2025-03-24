package apperrors

import (
	"encoding/json"
	"net/http"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	// BadRequest is for validation or invalid input errors
	BadRequest ErrorType = "BAD_REQUEST"
	// NotFound is for resource not found errors
	NotFound ErrorType = "NOT_FOUND"
	// Unauthorized is for authentication errors
	Unauthorized ErrorType = "UNAUTHORIZED"
	// Forbidden is for authorization errors
	Forbidden ErrorType = "FORBIDDEN"
	// Conflict is for resource conflicts (e.g., duplicate email)
	Conflict ErrorType = "CONFLICT"
	// InternalServer is for server errors
	InternalServer ErrorType = "INTERNAL_SERVER_ERROR"
)

// AppError represents an application error
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"-"` // HTTP status code, not exposed in JSON
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Type:    BadRequest,
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFound,
		Message: message,
		Code:    http.StatusNotFound,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    Unauthorized,
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    Forbidden,
		Message: message,
		Code:    http.StatusForbidden,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    Conflict,
		Message: message,
		Code:    http.StatusConflict,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string) *AppError {
	return &AppError{
		Type:    InternalServer,
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

// WriteError writes an error response to the HTTP response writer
func WriteError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"type":    err.Type,
			"message": err.Message,
		},
	}

	json.NewEncoder(w).Encode(response)
}
