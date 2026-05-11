package worker

import (
	"fmt"
	"os"

	"github.com/hibiken/asynq"
)

func NewRedisClient() *asynq.Client {
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}
