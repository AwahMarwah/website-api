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
	// CommissionRateBP dalam basis points: 1000 = 10%.
	CommissionRateBP int
	CreatedAt        time.Time
	UpdatedAt        *time.Time
}