package order

import "time"

type (
	CheckoutResponse struct {
		OrderID      string     `json:"order_id"`
		PaymentToken string     `json:"payment_token"`
		PaymentURL   *string    `json:"payment_url"`
		ExpiresAt    *time.Time `json:"expires_at"`
	}

	NotificationResponse struct {
		OrderID   string `json:"order_id"`
		Status    string `json:"status"`
		Processed bool   `json:"processed"`
	}

	PaymentLinkResponse struct {
		OrderID string `json:"order_id"`
		URL     string `json:"url"`
	}
)
