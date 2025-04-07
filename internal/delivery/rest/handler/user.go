package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/domain/user"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/apperrors"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService usecase.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService usecase.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// godoc GetUser
// @Summary Get User
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} user.User "Get user response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 404 {object} apperrors.AppError "Not Found"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid user ID"))
		return
	}

	// Get the user
	u, err := h.userService.GetUser(r.Context(), dtos.GetUserInput{ID: &userID})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			apperrors.WriteError(w, apperrors.NewNotFoundError("User not found"))
		default:
			apperrors.WriteError(w, apperrors.NewInternalServerError("Failed to get user"))
		}
		return
	}

	// Return the user
	response := map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// godoc UpdateUser
// @Summary Update User
// @Description Update a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param updateUserInput body dtos.UpdateUserInput true "Update user input"
// @Success 200 {object} user.User "Update user response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 404 {object} apperrors.AppError "Not Found"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid user ID"))
		return
	}

	// Parse the request body
	var input struct {
		Name     *string `json:"name,omitempty"`
		Password *string `json:"password,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid request body"))
		return
	}

	// Update the user
	updateInput := dtos.UpdateUserInput{
		ID:       userID,
		Name:     input.Name,
		Password: input.Password,
	}

	u, err := h.userService.UpdateUser(r.Context(), updateInput)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			apperrors.WriteError(w, apperrors.NewNotFoundError("User not found"))
		default:
			apperrors.WriteError(w, apperrors.NewInternalServerError("Failed to update user"))
		}
		return
	}

	// Return the updated user
	response := map[string]interface{}{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"role":  u.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// godoc ListUsers
// @Summary List Users
// @Description List all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []user.User "List users response"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	offset := 0
	limit := 10

	// Get users
	users, err := h.userService.ListUsers(r.Context(), dtos.ListUsersInput{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		apperrors.WriteError(w, apperrors.NewInternalServerError("Failed to list users"))
		return
	}

	// Map users to response format
	var usersResponse []map[string]interface{}
	for _, u := range users {
		usersResponse = append(usersResponse, map[string]interface{}{
			"id":     u.ID,
			"email":  u.Email,
			"name":   u.Name,
			"role":   u.Role,
			"status": u.Status,
		})
	}

	// Return the users
	response := map[string]interface{}{
		"users": usersResponse,
		"meta": map[string]interface{}{
			"offset": offset,
			"limit":  limit,
			"total":  len(usersResponse),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
