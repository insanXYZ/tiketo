package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type MockRedisClient struct{}

func NewMockRedisClient() RedisInterface {
	return &MockRedisClient{}
}

func (r *MockRedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}

func (r *MockRedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, redis.Nil
}

func (r *MockRedisClient) Del(ctx context.Context, key ...string) {

}
