package worker

import (
	"fmt"
	"log"
	"os"
	"website-api/task"

	"github.com/hibiken/asynq"
)

func StartWorker() {
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(task.TypeSendResetPassword, HandleResetPassword)
	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}

}
