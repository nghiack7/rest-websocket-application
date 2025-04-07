package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/personal/task-management/docs"
	"github.com/personal/task-management/internal/delivery/rest/handler"
	"github.com/personal/task-management/internal/delivery/rest/middleware"
	"github.com/personal/task-management/internal/delivery/websocket"
	httpserver "github.com/personal/task-management/pkg/server/http-server"
	"github.com/personal/task-management/pkg/utils/jwt"
)

// ServerDependencies holds all dependencies required for the server.
type ServerDependencies struct {
	UserHandler      *handler.UserHandler
	TaskHandler      *handler.TaskHandler
	AuthHandler      *handler.AuthHandler
	ChatHandler      *handler.ChatHandler
	JWTService       jwt.JWTTokenServicer
	RBACService      middleware.CasbinRBACService
	WebSocketHandler *websocket.Handler
}

func NewHTTPServer(cfg *viper.Viper, userHandler *handler.UserHandler, taskHandler *handler.TaskHandler, authHandler *handler.AuthHandler, rbacService middleware.CasbinRBACService, wsHandler *websocket.Handler, chatHandler *handler.ChatHandler) *httpserver.Server {
	host := cfg.GetString("server.host")
	port := cfg.GetInt("server.port")

	jwtService := jwt.NewJWTTokenService(cfg)

	dependencies := &ServerDependencies{
		UserHandler:      userHandler,
		TaskHandler:      taskHandler,
		AuthHandler:      authHandler,
		ChatHandler:      chatHandler,
		JWTService:       jwtService,
		RBACService:      rbacService,
		WebSocketHandler: wsHandler,
	}

	r := SetupRoutes(dependencies)
	return httpserver.NewServer(r, httpserver.WithServerHost(host), httpserver.WithServerPort(port))
}

// SetupRoutes initializes all application routes.
func SetupRoutes(deps *ServerDependencies) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", healthCheck)
	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.HandleFunc("/ws", deps.WebSocketHandler.HandleWebSocket)

	r.Route("/api", func(r chi.Router) {
		authRoutes(r, deps)
		userRoutes(r, deps)
		taskRoutes(r, deps)
		chatRoutes(r, deps)
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

func chatRoutes(router chi.Router, deps *ServerDependencies) {
	router.Route("/chat", func(r chi.Router) {
		// Room management
		r.Post("/rooms/direct", applyMiddlewares(deps.ChatHandler.CreateDirectRoom, deps))
		r.Post("/rooms/group", applyMiddlewares(deps.ChatHandler.CreateGroupRoom, deps))
		r.Get("/rooms", applyMiddlewares(deps.ChatHandler.ListRooms, deps))
		r.Get("/rooms/{roomId}", applyMiddlewares(deps.ChatHandler.GetRoomHistory, deps))
		r.Post("/rooms/{roomId}/join", applyMiddlewares(deps.ChatHandler.JoinRoom, deps))
		r.Post("/rooms/{roomId}/leave", applyMiddlewares(deps.ChatHandler.LeaveRoom, deps))
		r.Put("/rooms/{roomId}", applyMiddlewares(deps.ChatHandler.UpdateRoom, deps))

		// Message management
		r.Get("/rooms/{roomId}/messages", applyMiddlewares(deps.ChatHandler.GetMessages, deps))
		r.Post("/rooms/{roomId}/messages", applyMiddlewares(deps.ChatHandler.SendMessage, deps))
		r.Post("/rooms/{roomId}/messages/{messageId}/read", applyMiddlewares(deps.ChatHandler.MarkMessageAsRead, deps))
		r.Post("/rooms/{roomId}/messages/{messageId}/pin", applyMiddlewares(deps.ChatHandler.PinMessage, deps))
		r.Delete("/rooms/{roomId}/messages/{messageId}/pin", applyMiddlewares(deps.ChatHandler.UnpinMessage, deps))

		// Room actions
		r.Post("/rooms/{roomId}/archive", applyMiddlewares(deps.ChatHandler.ArchiveRoom, deps))
		r.Post("/rooms/{roomId}/unarchive", applyMiddlewares(deps.ChatHandler.UnarchiveRoom, deps))
		r.Post("/rooms/{roomId}/mute", applyMiddlewares(deps.ChatHandler.MuteRoom, deps))
		r.Post("/rooms/{roomId}/unmute", applyMiddlewares(deps.ChatHandler.UnmuteRoom, deps))
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
