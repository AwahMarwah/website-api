package worker

import (
	"context"
	"encoding/json"
	"log"
	lib "website-api/library/helper/email"
	userModel "website-api/model/user"
	"website-api/task"

	"github.com/hibiken/asynq"
)

func HandleResetPassword(ctx context.Context, t *asynq.Task) error {
	var payload task.ResetPasswordPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	user := userModel.User{
		Name:  payload.Name,
		Email: payload.Email,
	}

	err := lib.SendResetPasswordByEmail(&user, payload.Token)
	if err != nil {
		return err
	}

	log.Printf("email sent to %s", payload.Email)
	
	return nil
}
