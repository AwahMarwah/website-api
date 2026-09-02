package worker

import (
	"context"
	"encoding/json"
	"log"
	lib "website-api/library/helper/email"
	"website-api/task"

	"github.com/hibiken/asynq"
)

func HandlePaymentSuccess(ctx context.Context, t *asynq.Task) error {
	var payload task.PaymentSuccessPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	err := lib.SendPaymentSuccessByEmail(
		payload.Name,
		payload.Email,
		payload.OrderID,
		payload.TotalAmount,
	)
	if err != nil {
		return err
	}

	log.Printf("payment success email sent to %s for order %s", payload.Email, payload.OrderID)

	return nil
}
