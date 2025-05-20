## backend part

### initialize project
run the command as follow:
```
go mod init gin_im
```
The structure of the project:
```
project-root/  
├── config/  
│   ├── config.go  
│   └── config.yaml  
├── model/  
│   └── user.go  
├── repository/  
│   └── user.go  
├── service/  
│   └── user.go  
├── router/  
│   └── router.go  
├── wire/  
│   ├── wire.go  
│   └── wire_gen.go  
├── handle/  
│   └── user.go  
├── db/  
│   ├── mysql.go  
│   └── redis.go  
├── middleware/  
│   ├── jwt.go  
│   └── cors.go  
├── main.go  
├── go.mod  
└── go.sum  
```

### wire
Install and run 'wire': 
```
go install github.com/google/wire/cmd/wire@latest           // wire
wire gen ./wire.go                                          // generate the wire_gen.go file
```

### swagger
1. Install swagger
```
go install github.com/swaggo/swag/cmd/swag@latest           // swagger
```
2. Add swagger comment
``` go
// @title User API
// @version 1.0
// @description This is user API.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api/user
package handler

import (
	"gin_im/api"
	"gin_im/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}


// @Summary Register
// @Description Register
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body api.RegisterRequest true "User registration information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/users [post]
func (handle *UserHandler) Register(ctx *gin.Context) {

}
``` 

3. Generate Swagger documents
```
swag init -g main.go
```

4. Add swagger router
```go
package router

import (
    "gin_im/docs" // 替换为你的项目路径

	"github.com/gin-gonic/gin"
	
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	ginServer := gin.Default()

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
		api.POST("/register", userHandler.Register)
		api.POST("login", userHandler.Login)

	}

	return ginServer
}
```

5. Update `main.go`
```go
package main

import (
    // ...
	_ "gin_im/docs" 		// 这行必须存在
)

func main() {


}
```

6. Visit
http://localhost:8080/swagger/index.html


They will be installed at 'GOPATH', you can run 'go env GOPATH' to get the 'GOPATH'

### debug
launch.json
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "D:/code/gin_im/backend/main.go"			// the entrypoint of main.go
        }
    ]
}
```