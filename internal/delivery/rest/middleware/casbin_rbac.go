package middleware

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/personal/task-management/internal/domain/user"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type CasbinRBACService interface {
	HasPermission(role user.Role, resource string, action string) bool
	ApplyResourceFilter(r *http.Request, role user.Role, userID uuid.UUID)
}

// CasbinRBACService handles role-based access control using Casbin
type casbinRBACService struct {
	enforcer *casbin.Enforcer
}

// NewCasbinRBACService creates a new Casbin RBAC service
func NewCasbinRBACService(cfg *viper.Viper, db *gorm.DB) (CasbinRBACService, error) {
	enforcer, err := newCasbinEnforcer(cfg, db)
	if err != nil {
		return nil, err
	}
	// add policy to the enforcer
	enforcer.AddPolicy("employer", "tasks", "create")
	enforcer.AddPolicy("employer", "tasks", "read")
	enforcer.AddPolicy("employer", "tasks", "update")
	enforcer.AddPolicy("employer", "tasks", "delete")
	enforcer.AddPolicy("employer", "users", "create")
	enforcer.AddPolicy("employer", "users", "read")
	enforcer.AddPolicy("employer", "users", "update")
	enforcer.AddPolicy("employer", "users", "delete")
	enforcer.AddPolicy("employee", "tasks", "read")
	enforcer.AddPolicy("employee", "tasks", "update")
	enforcer.AddPolicy("employee", "users", "read")
	service := &casbinRBACService{
		enforcer: enforcer,
	}

	return service, nil
}

// HasPermission checks if a user has permission to perform an action on a resource
func (s *casbinRBACService) HasPermission(role user.Role, resource string, action string) bool {
	// Convert role to string
	roleStr := role.String()

	// Check permission using Casbin
	ok, err := s.enforcer.Enforce(roleStr, resource, action)
	if err != nil {
		return false
	}
	return ok
}

// ApplyResourceFilter applies resource filtering based on user role and permissions
func (s *casbinRBACService) ApplyResourceFilter(r *http.Request, role user.Role, userID uuid.UUID) {
	// Get the resource from path
	resource := GetResourceFromPath(r.URL.Path)
	if resource == "" {
		return
	}

	// For employees, add their user ID as a filter for their own resources
	if role == user.Employee {
		q := r.URL.Query()
		q.Set("assignee_id", userID.String())
		r.URL.RawQuery = q.Encode()
	}
}

// GetResourceFromPath extracts the resource from the request path
func GetResourceFromPath(path string) string {
	if strings.HasPrefix(path, "/api/tasks") {
		return "tasks"
	}
	if strings.HasPrefix(path, "/api/users") {
		return "users"
	}
	return ""
}

// GetActionFromMethod converts HTTP method to action
func GetActionFromMethod(method string) string {
	switch method {
	case http.MethodPost:
		return "create"
	case http.MethodGet:
		return "read"
	case http.MethodPut, http.MethodPatch:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return ""
	}
}
