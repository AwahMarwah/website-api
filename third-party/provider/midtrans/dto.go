package midtrans

import (
	"github.com/midtrans/midtrans-go"
)

// PaymentLinkItemDetail detail item untuk payment link
type PaymentLinkItemDetail struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int32  `json:"quantity"`
}

// PaymentLinkCustomerDetail data customer untuk payment link
type PaymentLinkCustomerDetail struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

// PaymentLinkRequest request untuk membuat payment link Midtrans
// https://docs.midtrans.com/reference/create-payment-link
type PaymentLinkRequest struct {
	TransactionDetails midtrans.TransactionDetails `json:"transaction_details"`
	Items              []PaymentLinkItemDetail     `json:"items,omitempty"`
	CustomerDetails    *PaymentLinkCustomerDetail  `json:"customer_details,omitempty"`
	UsageLimit         int                         `json:"usage_limit,omitempty"`
	Expiry             *ExpiryUnit                 `json:"expiry,omitempty"`
	EnabledPayments    []string                    `json:"enabled_payments,omitempty"`
}

// ExpiryUnit durasi kedaluwarsa payment link
type ExpiryUnit struct {
	StartTime string `json:"start_time,omitempty"`
	Unit      string `json:"unit"`
	Duration  int    `json:"duration"`
}

// PaymentLinkResponse response pembuatan payment link
type PaymentLinkResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	PaymentID string `json:"payment_id"`
}
