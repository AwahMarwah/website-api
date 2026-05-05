package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Helper struct {
	client *redis.Client
}

func NewHelper(client *redis.Client) *Helper {
	return &Helper{client: client}
}

// Generate hashed key (safe for filter & pagination)
func GenerateKey(prefix string, payload interface{}) string {
	b, _ := json.Marshal(payload)
	hash := md5.Sum(b)
	return prefix + ":" + hex.EncodeToString(hash[:])
}

// Reusable read-through cache
func (h *Helper) Cacheable(ctx context.Context, key string, ttl time.Duration, dest interface{}, fetch func() error) error {

	//  Try cache
	val, err := h.client.Get(ctx, key).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(val), dest); err == nil {
			return nil
		}
	}

	//  Fetch fresh data
	if err := fetch(); err != nil {
		return err
	}

	//  Store cache
	b, err := json.Marshal(dest)
	if err != nil {
		return nil
	}

	_ = h.client.Set(ctx, key, b, ttl).Err()

	return nil
}
