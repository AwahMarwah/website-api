package worker

import (
	"context"
	"encoding/json"
	"website-api/cache"
	"website-api/queue"

	"github.com/hibiken/asynq"
)

type CacheWorker struct {
	cache cache.Cache
}

func NewCacheWorker(cache cache.Cache) *CacheWorker {
	return &CacheWorker{cache: cache}
}

func (w *CacheWorker) HandleInvalidateCache(
	ctx context.Context,
	t *asynq.Task,
) error {

	var payload queue.InvalidateCachePayload

	if err := json.Unmarshal(
		t.Payload(),
		&payload,
	); err != nil {
		return err
	}

	for _, key := range payload.Keys {

		if err := w.cache.Delete(
			key,
		); err != nil {
			return err
		}
	}

	return nil
}
