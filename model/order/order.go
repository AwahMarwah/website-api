package order

import "time"

type Order struct {
	ID            string
	UserID        string
	AddressID     string
	TotalAmount   float64
	ShippingFee   float64
	Status        string
	PaymentMethod string
	PaymentToken  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID               string
	OrderID          string
	ProductVariantID string
	Price            float64
	Qty              int
	Subtotal         float64
}
