package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/personal/task-management/internal/domain/user"
	"github.com/personal/task-management/pkg/apperrors"
	"github.com/personal/task-management/pkg/utils/jwt"
)

func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	// Apply middleware in reverse order (like defer)
	for i := len(mid) - 1; i >= 0; i-- {
		handler = mid[i](handler)
	}
	return handler
}

func AuthMiddleware(jwtService jwt.JWTTokenServicer) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// bearer token
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// validate token
			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				apperrors.WriteError(w, apperrors.NewUnauthorizedError("Invalid token"))
				return
			}

			// verify token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				apperrors.WriteError(w, apperrors.NewUnauthorizedError("Invalid token"))
				return
			}
			// set claims to request
			ctx := context.WithValue(r.Context(), "user", claims)
			r = r.WithContext(ctx)
			// call next handler
			next.ServeHTTP(w, r)
		})
	}
}

// AuthorizationMiddleware enforces role-based access control using Casbin
func AuthorizationMiddleware(jwtService jwt.JWTTokenServicer, rbacService CasbinRBACService) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user").(*jwt.UserClaims)
			if !ok {
				apperrors.WriteError(w, apperrors.NewUnauthorizedError("Invalid claims"))
				return
			}
			// Convert role string to user.Role
			var userRole user.Role
			switch claims.Role {
			case "employee":
				userRole = user.Employee
			case "employer":
				userRole = user.Employer
			default:
				apperrors.WriteError(w, apperrors.NewUnauthorizedError("Invalid role"))
				return
			}

			// Get resource and action from request
			resource := GetResourceFromPath(r.URL.Path)
			action := GetActionFromMethod(r.Method)

			if resource == "" || action == "" {
				apperrors.WriteError(w, apperrors.NewForbiddenError("Permission denied: invalid resource or action"))
				return
			}

			// Check permission using Casbin
			if !rbacService.HasPermission(userRole, resource, action) {
				apperrors.WriteError(w, apperrors.NewForbiddenError(fmt.Sprintf("Permission denied: %s %s", action, resource)))
				return
			}

			// Apply resource filtering based on role
			rbacService.ApplyResourceFilter(r, userRole, claims.UserID)
			// store userID in context
			next.ServeHTTP(w, r)
		})
	}
}
