package db

import (
	"context"
	"fmt"
	"gin_im/config"
	"time"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB, 						// 要连接的 Redis 数据库编号（Redis 默认有 16 个数据库，编号从 0 到 15）
		PoolSize: cfg.Redis.PoolSize,
	})

	// 创建一个带有超时时间的上下文 ctx，超时时间设置为 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// 使用 Redis.Ping(ctx) 方法发送一个 PING 命令到 Redis 服务器，以测试连接是否成功
	// 通过调用 .Result() 方法获取 PING 命令的响应结果。如果连接成功，PING 命令通常会返回一个简单的 "PONG" 响应。
	_, err := client.Ping(ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return client, nil
}