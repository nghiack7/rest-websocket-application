package mocks

//go:generate mockgen -destination=./user_repository.go -package=mocks github.com/personal/task-management/internal/repositories UserRepository
//go:generate mockgen -destination=./hasher.go -package=mocks github.com/personal/task-management/internal/usecase Hasher
//go:generate mockgen -destination=./jwt_service.go -package=mocks github.com/personal/task-management/pkg/utils/jwt JWTTokenServicer
//go:generate mockgen -destination=./user_service.go -package=mocks github.com/personal/task-management/internal/usecase UserService
//go:generate mockgen -destination=./task_service.go -package=mocks github.com/personal/task-management/internal/usecase TaskService
//go:generate mockgen -destination=./casbin_rbac_service.go -package=mocks github.com/personal/task-management/internal/delivery/rest/middleware CasbinRBACService
//go:generate mockgen -destination=./task_repository.go -package=mocks github.com/personal/task-management/internal/repositories TaskRepository
//go:generate mockgen -destination=./websocket_service.go -package=mocks github.com/personal/task-management/internal/usecase WebSocketService
