package user

import (
	"time"

	"github.com/google/uuid"
)

// Role represents user roles in the system
type Role int

const (
	// RoleGuest represents a guest role
	Unknown Role = iota
	// RoleEmployee represents an employee role
	Employee
	// RoleEmployer represents an employer role
	Employer
)

func (r Role) String() string {
	switch r {
	case Unknown:
		return "unknown"
	case Employee:
		return "employee"
	case Employer:
		return "employer"
	default:
		return "unknown"
	}
}

type Status int

const (
	StatusActive Status = iota
	StatusInactive
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	}
	return "unknown"
}

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // Never expose password
	Role      Role      `json:"role"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user with the given parameters
func NewUser(email, name, password string) (*User, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if password == "" {
		return nil, ErrEmptyPassword
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		Password:  password, // Note: Should be hashed before storage
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) SetRole(role string) {
	switch role {
	case "employee":
		u.Role = Employee
	case "employer":
		u.Role = Employer
	default:
		u.Role = Unknown
	}
}

// IsEmployer checks if user has employer role
func (u *User) IsEmployer() bool {
	return u.Role == Employer
}

// IsEmployee checks if user has employee role
func (u *User) IsEmployee() bool {
	return u.Role == Employee
}

// CanCreateTasks checks if user can create tasks
func (u *User) CanCreateTasks() bool {
	return u.IsEmployer()
}

// CanAssignTasks checks if user can assign tasks
func (u *User) CanAssignTasks() bool {
	return u.IsEmployer()
}

// CanViewAllTasks checks if user can view all tasks
func (u *User) CanViewAllTasks() bool {
	return u.IsEmployer()
}

// CanUpdateTaskStatus checks if user can update task status
func (u *User) CanUpdateTaskStatus() bool {
	return true // Both roles can update status, but employee only their assigned tasks
}
