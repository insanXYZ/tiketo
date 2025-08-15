package db

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, any, time.Duration) error
	Del(context.Context, ...string)
}

const (
	ExpRedis = time.Duration(5 * time.Minute)
)

type RedisClient struct {
	redis *redis.Client
}

func NewRedisClient() (RedisInterface, error) {
	envDb := os.Getenv("REDIS_DB")

	rdDb, err := strconv.Atoi(envDb)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       rdDb,
	})

	redisClient := &RedisClient{client}

	return redisClient, redisClient.redis.Ping(context.Background()).Err()

}

func (r *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return r.redis.Get(ctx, key).Bytes()
}

func (r *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.redis.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) {
	r.redis.Del(ctx, keys...)
}
