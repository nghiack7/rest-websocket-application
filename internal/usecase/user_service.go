package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/personal/task-management/internal/delivery/rest/dtos"
	"github.com/personal/task-management/internal/domain/user"
	repository "github.com/personal/task-management/internal/repositories"
	"github.com/personal/task-management/pkg/utils/jwt"
)

type UserService interface {
	RegisterUser(ctx context.Context, input dtos.RegisterUserInput) (*dtos.GetUserOutput, error)
	Login(ctx context.Context, input dtos.LoginInput) (*dtos.LoginOutput, error)
	GetUser(ctx context.Context, input dtos.GetUserInput) (*user.User, error)
	UpdateUser(ctx context.Context, input dtos.UpdateUserInput) (*user.User, error)
	ListUsers(ctx context.Context, input dtos.ListUsersInput) ([]*user.User, error)
}

// ErrInvalidCredentials is returned when authentication fails
var ErrInvalidCredentials = errors.New("invalid email or password")

// UserService handles user-related operations and business logic
type userService struct {
	userRepo     repository.UserRepository
	hasher       Hasher
	tokenService jwt.JWTTokenServicer
}

type Hasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, plainPassword string) bool
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository, hasher Hasher, tokenService jwt.JWTTokenServicer) UserService {
	return &userService{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

// RegisterUser registers a new user
func (s *userService) RegisterUser(ctx context.Context, input dtos.RegisterUserInput) (*dtos.GetUserOutput, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, user.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := s.hasher.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	newUser, err := user.NewUser(
		input.Email,
		input.Name,
		hashedPassword,
	)
	newUser.SetRole(input.Role)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}
	resp := &dtos.GetUserOutput{
		ID:     newUser.ID,
		Email:  newUser.Email,
		Name:   newUser.Name,
		Role:   newUser.Role.String(),
		Status: newUser.Status.String(),
	}

	return resp, nil
}

// Login authenticates a user and returns an auth token
func (s *userService) Login(ctx context.Context, input dtos.LoginInput) (*dtos.LoginOutput, error) {
	// Find user by email
	u, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if !s.hasher.ComparePasswords(u.Password, input.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenService.GenerateToken(u.ID, u.Email, u.Role.String())
	if err != nil {
		return nil, err
	}

	return &dtos.LoginOutput{
		User: &dtos.GetUserOutput{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role.String(),
		},
		AuthToken: token,
	}, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, input dtos.GetUserInput) (*user.User, error) {
	return s.userRepo.GetByID(ctx, *input.ID)
}

// UpdateUser updates a user's information
func (s *userService) UpdateUser(ctx context.Context, input dtos.UpdateUserInput) (*user.User, error) {
	// Get user
	u, err := s.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if input.Name != nil {
		u.Name = *input.Name
	}

	if input.Password != nil {
		hashedPassword, err := s.hasher.HashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		u.Password = hashedPassword
	}

	// Update timestamp
	u.UpdatedAt = time.Now()

	// Save user
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) ListUsers(ctx context.Context, input dtos.ListUsersInput) ([]*user.User, error) {
	return s.userRepo.List(ctx, input.Offset, input.Limit)
}
