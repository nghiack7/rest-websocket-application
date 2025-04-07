//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/personal/task-management/config"
	api "github.com/personal/task-management/internal/delivery/rest/handler"
	"github.com/personal/task-management/internal/delivery/rest/middleware"
	"github.com/personal/task-management/internal/delivery/websocket"
	"github.com/personal/task-management/internal/repositories/postgres"
	internalServer "github.com/personal/task-management/internal/server"
	"github.com/personal/task-management/internal/usecase"
	"github.com/personal/task-management/pkg/app"
	"github.com/personal/task-management/pkg/db"
	"github.com/personal/task-management/pkg/server/http-server"
	"github.com/personal/task-management/pkg/utils/hasher"
	"github.com/personal/task-management/pkg/utils/jwt"
)

func NewWire() (*app.App, func(), error) {
	panic(wire.Build(
		config.LoadConfig,
		db.ConnectDB,
		loadGormDB,
		postgres.NewPostgresUserRepository,
		postgres.NewPostgresTaskRepository,
		postgres.NewChatRepository,
		loadHasher,
		jwt.NewJWTTokenService,
		usecase.NewUserService,
		usecase.NewTaskService,
		usecase.NewWebSocketService,
		api.NewUserHandler,
		api.NewTaskHandler,
		api.NewAuthHandler,
		api.NewChatHandler,
		websocket.NewHandler,
		middleware.NewCasbinRBACService,
		internalServer.NewHTTPServer,
		newApp,
	))
}

func newApp(httpServer *http.Server) (*app.App, func(), error) {
	app := app.NewApp(app.WithServer(httpServer), app.WithName("task-management"))
	return app, func() {
		app.Stop()
	}, nil
}

func loadGormDB(instance *db.PostgresDB) *gorm.DB {
	instance.MigrateDB()
	return instance.GetDB()
}

func loadHasher(cfg *viper.Viper) usecase.Hasher {
	return hasher.NewBcryptHasher(cfg)
}
