package router

import (
	"gin_im/config"
	"gin_im/docs"
	"gin_im/handler"
	"gin_im/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	config *config.Config,
	userHandler *handler.UserHandler,
	websocketHandler *handler.WebsocketHandler,
) *gin.Engine {
	ginServer := gin.Default()

	ginServer.Use(middleware.Cors())

	// 程序化设置 Swagger 信息
	docs.SwaggerInfo.Title = "GIN IM System"
	docs.SwaggerInfo.Description = "GIN IM System API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 添加 Swagger 路由
	ginServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := ginServer.Group("/api")
	{
		user := api.Group("/user")
		{
			// 不需要认证的路由
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)

			// 需要认证的路由组
			authUser := user.Group("")
			authUser.Use(middleware.JWTAuth(config))
			{
				authUser.GET("/information", userHandler.GetUserInformation)
				authUser.PUT("/update", userHandler.UpdateUser)

			}
		}

		// WebSocket 相关路由
		ws := api.Group("/ws")
		ws.Use()
		// ws.Use(middleware.JWTAuth(config))
		{
			ws.GET("/connect", websocketHandler.Connect)
		}

	}

	return ginServer
}
