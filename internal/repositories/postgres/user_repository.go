package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/personal/task-management/internal/domain/user"
	repository "github.com/personal/task-management/internal/repositories"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *user.User) error {
	return r.db.Create(user).Error
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *user.User) error {
	return r.db.Save(user).Error
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Delete(&user.User{}, "id = ?", id).Error
}

func (r *PostgresUserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var users []*user.User
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
