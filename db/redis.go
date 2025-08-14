package db

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ExpRedis = time.Duration(5 * time.Minute)
)

func NewRedisClient() (*redis.Client, error) {
	envDb := os.Getenv("REDIS_DB")

	rdDb, err := strconv.Atoi(envDb)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       rdDb,
	})

	return redisClient, redisClient.Ping(context.Background()).Err()
}
