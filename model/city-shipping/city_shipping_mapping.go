package city_shipping

import "time"

type CityShippingMapping struct {
	CityID        uint64    `json:"city_id"`
	Provider      string    `json:"provider"`
	DestinationId uint64    `json:"destination_id"`
	CreatedAt     time.Time `json:"created_at"`
}
