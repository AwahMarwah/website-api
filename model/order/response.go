package order

import "time"

type (
	CheckoutResponse struct {
		OrderID      string     `json:"order_id"`
		PaymentToken string     `json:"payment_token"`
		PaymentURL   *string    `json:"payment_url"`
		ExpiresAt    *time.Time `json:"expires_at"`
		TotalAmount  float64    `json:"total_amount"`
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
		ProductID        string  `json:"product_id"`
		ProductName      string  `json:"product_name"`
		VariantName      string  `json:"variant_name"`
		ProductImageURL  string  `json:"product_image_url"`
		Sku              string  `json:"sku"`
		MerchantID       string  `json:"merchant_id"`
		MerchantName     string  `json:"merchant_name"`
		Reviewed         bool    `json:"reviewed"`
		ProductVariantID string  `json:"product_variant_id"`
		Price            float64 `json:"price"`
		Qty              int     `json:"qty"`
		Subtotal         float64 `json:"subtotal"`
	}

	OrderStatusHistoryResponse struct {
		FromStatus string    `json:"from_status"`
		ToStatus   string    `json:"to_status"`
		ActorRole  string    `json:"actor_role"`
		Note       *string   `json:"note"`
		CreatedAt  time.Time `json:"created_at"`
	}

	OrderMerchantShippingResponse struct {
		MerchantID   string `json:"merchant_id"`
		MerchantName string `json:"merchant_name"`
		Courier      string `json:"courier"`
		Service      string `json:"service"`
		Cost         int64  `json:"cost"`
		Etd          string `json:"etd"`
		WeightGram   int    `json:"weight_gram"`
	}

	OrderResponse struct {
		ID              string                          `json:"id"`
		AddressID       string                          `json:"address_id"`
		AddressSnapshot AddressSnapshot                 `json:"address_snapshot"`
		TotalAmount     float64                         `json:"total_amount"`
		ShippingFee     float64                         `json:"shipping_fee"`
		DiscountAmount  float64                         `json:"discount_amount"`
		Status          string                          `json:"status"`
		PaymentMethod   string                          `json:"payment_method"`
		PaymentURL      *string                         `json:"payment_url"`
		ExpiredAt       *time.Time                      `json:"expired_at"`
		Note            *string                         `json:"note"`
		PaidAt          *time.Time                      `json:"paid_at"`
		ShippedAt       *time.Time                      `json:"shipped_at"`
		CompletedAt     *time.Time                      `json:"completed_at"`
		CancelledAt     *time.Time                      `json:"cancelled_at"`
		CreatedAt       time.Time                       `json:"created_at"`
		Items           []OrderItemResponse             `json:"items"`
		StatusHistories []OrderStatusHistoryResponse    `json:"status_histories"`
		Shippings       []OrderMerchantShippingResponse `json:"shippings"`
	}
)
