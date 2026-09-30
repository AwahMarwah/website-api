package voucher

import "time"

type Voucher struct {
	ID           string
	Code         string
	Description  *string
	Type         string
	Value        float64
	MaxDiscount  *float64
	MinSpend     float64
	Quota        *int
	UsedCount    int
	PerUserLimit int
	StartsAt     *time.Time
	EndsAt       *time.Time
	MerchantID   *string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type ProductVoucher struct {
	VoucherID string
	ProductID string
}

type VoucherRedemption struct {
	ID              string
	VoucherID       string
	OrderID         string
	UserID          string
	DiscountAmount  float64
	CreatedAt       time.Time
}
