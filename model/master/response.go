package master

type (
	// ===== PROVINCE ===== //
	ListProvinceResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// ===== DISTRICT ===== //
	ListDistrictResponse struct {
		ID     string `json:"id"`
		CityID string `json:"city_id"`
		Name   string `json:"name"`
	}

	// ===== SUBDISTRICT ===== //
	ListSubdistrictResponse struct {
		ID         string `json:"id"`
		DistrictID string `json:"district_id"`
		Name       string `json:"name"`
	}

	// ===== CITY ===== //
	ListCityResponse struct {
		ID         string `json:"id"`
		ProvinceID string `json:"province_id"`
		Name       string `json:"name"`
	}

	// ===== SWAGGER DOCS ===== //
	SwaggerProvincePagination struct {
		Data    []ListProvinceResponse `json:"data"`
		Message string                 `json:"message" example:"OK"`
		Page    struct {
			Current int `json:"current" example:"1"`
			Size    int `json:"size" example:"10"`
			Total   int `json:"total" example:"37"`
		} `json:"page"`
	}
	SwaggerDistrictPagination struct {
		Data    []ListDistrictResponse `json:"data"`
		Message string                 `json:"message" example:"OK"`
		Page    struct {
			Current int `json:"current" example:"1"`
			Size    int `json:"size" example:"10"`
			Total   int `json:"total" example:"37"`
		} `json:"page"`
	}

	SwaggerCityPagination struct {
		Data    []ListCityResponse `json:"data"`
		Message string             `json:"message" example:"OK"`
		Page    struct {
			Current int `json:"current" example:"1"`
			Size    int `json:"size" example:"10"`
			Total   int `json:"total" example:"37"`
		} `json:"page"`
	}

	SwaggerSubdistrictPagination struct {
		Data    []ListSubdistrictResponse `json:"data"`
		Message string                    `json:"message" example:"OK"`
		Page    struct {
			Current int `json:"current" example:"1"`
			Size    int `json:"size" example:"10"`
			Total   int `json:"total" example:"37"`
		} `json:"page"`
	}
)
