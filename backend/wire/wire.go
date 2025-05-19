//go:build wireinject
// +build wireinject

package wire

import (
	"gin_im/db"
	"gin_im/router"
	"gin_im/config"
	"gin_im/handler"
	"gin_im/service"
	"gin_im/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var handlerSet = wire.NewSet(
	handler.NewUserHandler,
)

var serviceSet = wire.NewSet(
	service.NewUserService,
)

var repositorySet =wire.NewSet(
	repository.NewUserRepository,
)

var dbSet = wire.NewSet(
	db.NewMySQLDB,
	db.NewRedisClient,
)

var routerSet = wire.NewSet(
	router.SetupRouter,
)

var configSet = wire.NewSet(
	config.LoadConfig,
)

var SuperSet = wire.NewSet(
	handlerSet,
	serviceSet,
	repositorySet,
	dbSet,
	routerSet,
	configSet,
)

func InitializeApp(configPath string) (*gin.Engine, error) {
	wire.Build(SuperSet)

	return &gin.Engine{}, nil
}
