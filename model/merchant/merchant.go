package merchant

import "time"

type Merchant struct {
	ID            string
	Name          string
	Slug          string
	DestinationID int64
	CityID        string
	Address       string
	IsActive      bool
	UserID        string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}