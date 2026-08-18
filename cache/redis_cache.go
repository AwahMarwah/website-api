package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return ErrCacheMiss
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) DeleteByPattern(pattern string) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	iter := r.client.Scan(
		ctx,
		0,
		pattern,
		0,
	).Iterator()

	for iter.Next(ctx) {

		if err := r.client.Del(
			ctx,
			iter.Val(),
		).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}
