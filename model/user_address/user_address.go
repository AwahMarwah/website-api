package user_address

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAddress struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	RecipientName string    `json:"recipient_name"`
	PhoneNumber   string    `json:"phone_number"`
	FullAddress   string    `json:"full_address"`
	City          string    `json:"city"`
	PostalCode    string    `json:"postal_code"`
	IsPrimary     bool      `json:"is_primary"`
	ProvinceID    string    `json:"province_id"`
	CityID        string    `json:"city_id"`
	DistrictID    string    `json:"district_id"`
	SubdistrictID string    `json:"subdistrict_id"`
	DestinationID int64     `json:"destination_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (userAddress *UserAddress) BeforeCreate(*gorm.DB) error {
	userAddress.ID = uuid.New().String()
	return nil
}
