package master

import "time"

type (
	Province struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	City struct {
		ID         string `json:"id"`
		ProvinceID string `json:"province_id"`
		Code       string `json:"code"`
		Name       string `json:"name"`
	}

	District struct {
		ID     string `json:"id"`
		CityID string `json:"city_id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
	}

	Subdistrict struct {
		ID         string    `json:"id"`
		DistrictID string    `json:"district_id"`
		Name       string    `json:"name"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
)
