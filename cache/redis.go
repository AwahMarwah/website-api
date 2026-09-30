package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	addr := host + ":" + port
	log.Printf("redis: connecting to %s (db %d)", addr, db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}

	log.Printf("redis: connected to %s", addr)
	return rdb
}
