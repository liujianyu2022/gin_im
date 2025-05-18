package main

import (
	
	"log"

	"github.com/gin-gonic/gin"

	"gin_im/config"
	"gin_im/wire"
)

func main() {
	// 加载配置
	configPath := "config/config.yaml"
	cfg := config.LoadConfig(configPath)

	// 设置Gin模式
	gin.SetMode(cfg.App.Mode)

	// 初始化依赖
	router, err := wire.InitializeApp(configPath)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// 启动服务
	if err := router.Run(cfg.App.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}