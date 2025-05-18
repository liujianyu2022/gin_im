//go:build wireinject
// +build wireinject

package wire

import (
	"gin_im/db"
	"gin_im/router"
	"gin_im/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
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
	dbSet,
	routerSet,
)

func InitializeApp(configPath string) (*gin.Engine, error) {
	wire.Build(SuperSet)

	return &gin.Engine{}, nil
}
