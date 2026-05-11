package task

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSendResetPassword = "email:reset_password"
	TypeSendVerification  = "email:verification"
)

type ResetPasswordPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewResetPasswordTask(name, email, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(ResetPasswordPayload{
		Name:  name,
		Email: email,
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendResetPassword, payload), nil

}
