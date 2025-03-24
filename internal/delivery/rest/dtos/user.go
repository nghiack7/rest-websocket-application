package dtos

import "github.com/google/uuid"

type RegisterUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=employee employer"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginOutput struct {
	User      *GetUserOutput `json:"user"`
	AuthToken string         `json:"auth_token"`
}

type GetUserInput struct {
	ID    *uuid.UUID `json:"id"`
	Email *string    `json:"email""`
}

type UpdateUserInput struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Name     *string   `json:"name"`
	Password *string   `json:"password"`
}

type ListUsersInput struct {
	Offset int    `json:"offset" validate:"min=0"`
	Limit  int    `json:"limit" validate:"required,min=1,max=100"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	SortBy string `json:"sort_by" validate:"oneof=name email role"`
	Role   string `json:"role" validate:"oneof=employee employer"`
	Status string `json:"status" validate:"oneof=active inactive"`
	Search string `json:"search"`
}

type GetUserOutput struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Status string    `json:"status"`
}
