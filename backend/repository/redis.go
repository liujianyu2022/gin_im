package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)


type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) Publish(ctx *gin.Context, channel string, message any) error {
	_, err := r.client.Publish(ctx, channel, message).Result()
	return err
}

func (r *RedisRepository) Subscribe(ctx *gin.Context, channels ...string) (string, error) {
	pubsub := r.client.Subscribe(ctx, channels...)
	defer pubsub.Close()

	msg, err := pubsub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}

	return msg.Payload, nil
}