package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/domain/user"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/apperrors"
)

type AuthHandler struct {
	userService usecase.UserService
}

func NewAuthHandler(userService usecase.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// godoc Login
// @Summary Login
// @Description Login to the system
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInput body dtos.LoginInput true "Login input"
// @Success 200 {object} map[string]interface{} "Login response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 401 {object} apperrors.AppError "Unauthorized"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var input dtos.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid request body"))
		return
	}

	// Authenticate the user
	authUser, err := h.userService.Login(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCredentials):
			apperrors.WriteError(w, apperrors.NewUnauthorizedError("Invalid email or password"))
		default:
			apperrors.WriteError(w, apperrors.NewInternalServerError("Failed to login"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authUser)
}

// godoc RegisterUser
// @Summary Register User
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param registerUserInput body dtos.RegisterUserInput true "Register user input"
// @Success 201 {object} dtos.GetUserOutput "Register response"
// @Failure 400 {object} apperrors.AppError "Bad Request"
// @Failure 409 {object} apperrors.AppError "Conflict"
// @Failure 500 {object} apperrors.AppError "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var input dtos.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.WriteError(w, apperrors.NewBadRequestError("Invalid request body"))
		return
	}

	// Register the user
	newUser, err := h.userService.RegisterUser(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailExists):
			apperrors.WriteError(w, apperrors.NewConflictError("Email already exists"))
		default:
			apperrors.WriteError(w, apperrors.NewInternalServerError("Failed to register user"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
