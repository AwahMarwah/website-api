package order

import (
	"time"

	"website-api/model/product-variant"
)

type Order struct {
	ID            string
	UserID        string
	AddressID     string
	TotalAmount   float64
	ShippingFee   float64
	Status        string
	PaymentMethod string
	PaymentToken  string
	PaymentURL    *string
	ExpiredAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID               string
	OrderID          string
	ProductVariantID string
	Price            float64
	Qty              int
	Subtotal         float64
	TotalWeightGram  int
	ProductVariant   *product_variant.ProductVariant `gorm:"foreignKey:ProductVariantID"`
}

type OrderMerchantShipping struct {
	ID          string
	OrderID     string
	MerchantID  string
	Courier     string
	Service     string
	Cost        int64
	Etd         string
	WeightGram  int
	CreatedAt   time.Time
}
