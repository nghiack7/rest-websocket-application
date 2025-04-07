package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/personal/task-management/internal/domain/user"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create stores a new user in the repository
	Create(ctx context.Context, user *user.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*user.User, error)

	// Update updates an existing user in the repository
	Update(ctx context.Context, user *user.User) error

	// Delete removes a user from the repository
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with optional pagination
	List(ctx context.Context, offset, limit int) ([]*user.User, error)
}
