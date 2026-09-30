package voucher

import "time"

type VoucherResponse struct {
	ID            string     `json:"id"`
	Code          string     `json:"code"`
	Description   *string    `json:"description"`
	Type          string     `json:"type"`
	Value         float64    `json:"value"`
	MaxDiscount   *float64   `json:"max_discount"`
	MinSpend      float64    `json:"min_spend"`
	Quota         *int       `json:"quota"`
	UsedCount     int        `json:"used_count"`
	PerUserLimit  int        `json:"per_user_limit"`
	RemainingQuota *int      `json:"remaining_quota"`
	StartsAt      *time.Time `json:"starts_at"`
	EndsAt        *time.Time `json:"ends_at"`
	MerchantID    *string    `json:"merchant_id"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ValidateVoucherResponse struct {
	Code            string  `json:"code"`
	Valid           bool    `json:"valid"`
	Message         string  `json:"message"`
	DiscountAmount  float64 `json:"discount_amount"`
	FinalAmount     float64 `json:"final_amount"`
}
