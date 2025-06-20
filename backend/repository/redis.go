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


// 在 repository 层添加获取群组成员的方法
func (repository *RedisRepository) GetGroupMembers(groupId int64) ([]int64, error) {
	// 实际实现可以从Redis或数据库获取

	// 这里用模拟数据示例
	return []int64{1001, 1002, 1003}, nil
}