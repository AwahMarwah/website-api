package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func OpenRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		PoolSize: 20,
	})

	return rdb, nil
}
