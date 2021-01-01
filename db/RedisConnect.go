package db

import (
	"context"
	"davidws/config"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConnect connects to a Redis instance and returns a *redis.Client and an error
func RedisConnect(ctx context.Context) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:        config.RedisAddr,
		Password:    config.RedisDBPass,
		DB:          0,
		MaxConnAge:  time.Hour,
		IdleTimeout: 4,
	})

	_, err = client.Ping(ctx).Result()

	return
}
