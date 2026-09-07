package merchant

type (
	MerchantResponse struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		DestinationID int64  `json:"destination_id"`
		CityID        string `json:"city_id"`
		Address       string `json:"address"`
		IsActive      bool   `json:"is_active"`
	}
)