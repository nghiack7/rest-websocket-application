package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/personal/task-management/docs"
	"github.com/personal/task-management/internal/delivery/rest/handler"
	"github.com/personal/task-management/internal/delivery/rest/middleware"
	httpserver "github.com/personal/task-management/pkg/server/http-server"
	"github.com/personal/task-management/pkg/utils/jwt"
)

// ServerDependencies holds all dependencies required for the server.
type ServerDependencies struct {
	UserHandler *handler.UserHandler
	TaskHandler *handler.TaskHandler
	AuthHandler *handler.AuthHandler
	JWTService  jwt.JWTTokenServicer
	RBACService middleware.CasbinRBACService
}

func NewHTTPServer(cfg *viper.Viper, userHandler *handler.UserHandler, taskHandler *handler.TaskHandler, authHandler *handler.AuthHandler, rbacService middleware.CasbinRBACService) *httpserver.Server {
	host := cfg.GetString("server.host")
	port := cfg.GetInt("server.port")

	jwtService := jwt.NewJWTTokenService(cfg)

	dependencies := &ServerDependencies{
		UserHandler: userHandler,
		TaskHandler: taskHandler,
		AuthHandler: authHandler,
		JWTService:  jwtService,
		RBACService: rbacService,
	}

	r := SetupRoutes(dependencies)
	return httpserver.NewServer(r, httpserver.WithServerHost(host), httpserver.WithServerPort(port))
}

// SetupRoutes initializes all application routes.
func SetupRoutes(deps *ServerDependencies) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", healthCheck)
	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		authRoutes(r, deps)
		userRoutes(r, deps)
		taskRoutes(r, deps)
	})

	return r
}

func authRoutes(router chi.Router, deps *ServerDependencies) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", deps.AuthHandler.RegisterUser)
		r.Post("/login", deps.AuthHandler.Login)
	})
}

func userRoutes(router chi.Router, deps *ServerDependencies) {
	router.Route("/users", func(r chi.Router) {
		r.Get("/", applyMiddlewares(deps.UserHandler.ListUsers, deps))
		r.Get("/{id}", applyMiddlewares(deps.UserHandler.GetUser, deps))
		r.Put("/{id}", applyMiddlewares(deps.UserHandler.UpdateUser, deps))
	})
}

func taskRoutes(router chi.Router, deps *ServerDependencies) {
	router.Route("/tasks", func(r chi.Router) {
		r.Post("/", applyMiddlewares(deps.TaskHandler.Create, deps))
		r.Get("/", applyMiddlewares(deps.TaskHandler.List, deps))
		r.Get("/{id}", applyMiddlewares(deps.TaskHandler.Get, deps))
		r.Put("/{id}", applyMiddlewares(deps.TaskHandler.Update, deps))
		r.Delete("/{id}", applyMiddlewares(deps.TaskHandler.Delete, deps))
	})
}

// applyMiddlewares wraps a handler with authentication and authorization.
func applyMiddlewares(handlerFunc http.HandlerFunc, deps *ServerDependencies) http.HandlerFunc {
	return middleware.Use(handlerFunc,
		middleware.AuthMiddleware(deps.JWTService),
		middleware.AuthorizationMiddleware(deps.JWTService, deps.RBACService),
	)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
