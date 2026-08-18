package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
}

func New(client *asynq.Client) *Client {
	return &Client{client: client}
}

func (c *Client) InvalidateCache(
	keys []string,
) error {

	payload, _ := json.Marshal(
		InvalidateCachePayload{
			Keys: keys,
		},
	)

	task := asynq.NewTask(
		TaskInvalidateCache,
		payload,
	)

	_, err := c.client.Enqueue(task)

	return err
}
