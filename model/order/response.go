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

	OrderItemResponse struct {
		ID               string  `json:"id"`
		ProductVariantID string  `json:"product_variant_id"`
		Price            float64 `json:"price"`
		Qty              int     `json:"qty"`
		Subtotal         float64 `json:"subtotal"`
	}

	OrderResponse struct {
		ID            string             `json:"id"`
		AddressID     string             `json:"address_id"`
		TotalAmount   float64            `json:"total_amount"`
		ShippingFee   float64            `json:"shipping_fee"`
		Status        string             `json:"status"`
		PaymentMethod string             `json:"payment_method"`
		PaymentURL    *string            `json:"payment_url"`
		ExpiredAt     *time.Time         `json:"expired_at"`
		CreatedAt     time.Time          `json:"created_at"`
		Items         []OrderItemResponse `json:"items"`
	}
)
