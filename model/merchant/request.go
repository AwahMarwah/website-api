package merchant

type (
	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	CreateMerchantReq struct {
		Name          string `binding:"required" json:"name"`
		Slug          string `binding:"required" json:"slug"`
		DestinationID int64  `binding:"required" json:"destination_id"`
		CityID        string `json:"city_id"`
		Address       string `json:"address"`
	}

	UpdateMerchantReq struct {
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		DestinationID *int64 `json:"destination_id"`
		CityID        string `json:"city_id"`
		Address       string `json:"address"`
		IsActive      *bool  `json:"is_active"`
	}
)