package task

import (
	"github.com/hibiken/asynq"
)

const TypeCancelExpiredOrders = "order:cancel_expired"

func NewCancelExpiredOrdersTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeCancelExpiredOrders, nil), nil
}
