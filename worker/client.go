package worker

import (
	"fmt"
	"os"
	"sync"

	"github.com/hibiken/asynq"
)

var (
	client *asynq.Client
	once   sync.Once
)

func NewRedisClient() *asynq.Client {
	once.Do(func() {
		addr := fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		)

		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
		})
	})

	return client
}
