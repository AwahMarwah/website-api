package master

type (
	// ===== Province ===== //
	GetListProvinceRequest struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
		Search string `form:"search"`
	}

	// ===== CITY ===== //
	GetListCitiesRequest struct {
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Offset     int    `form:"offset"`
		Search     string `form:"search"`
		ProvinceID string `form:"province_id"`
	}

	// ===== DISTRICT ===== //
	GetListDistrictsRequest struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
		Search string `form:"search"`
		CityID string `form:"city_id"`
	}
)
