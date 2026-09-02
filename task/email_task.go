package task

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSendResetPassword = "email:reset_password"
	TypeSendVerification  = "email:verification"
	TypeSendPaymentSuccess = "email:payment_success"
)

type ResetPasswordPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type PaymentSuccessPayload struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	OrderID     string  `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
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

func NewPaymentSuccessTask(name, email, orderID string, totalAmount float64) (*asynq.Task, error) {
	payload, err := json.Marshal(PaymentSuccessPayload{
		Name:        name,
		Email:       email,
		OrderID:     orderID,
		TotalAmount: totalAmount,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeSendPaymentSuccess, payload), nil

}
