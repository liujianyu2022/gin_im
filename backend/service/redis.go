package service

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{client: client}
}

func (redis *RedisService) Publish(ctx *gin.Context, channel string, message any) error {
	_, err := redis.client.Publish(ctx, channel, message).Result()
	return err
}

func (redis *RedisService) Subscribe(ctx *gin.Context, channels ...string) *redis.PubSub {
	return redis.client.Subscribe(ctx, channels...)
}