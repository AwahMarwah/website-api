package user_address

import "time"

type UserAddress struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	RecipientName string    `json:"recipient_name"`
	PhoneNumber   string    `json:"phone_number"`
	FullAddress   string    `json:"full_address"`
	City          string    `json:"city"`
	PostalCode    string    `json:"postal_code"`
	IsPrimary     bool      `json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
