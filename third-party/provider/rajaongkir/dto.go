package rajaongkir

type (
	Meta struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
		Status  string `json:"status"`
	}

	// ==========================================
	// Destination
	// ==========================================

	DestinationResponse struct {
		Meta Meta          `json:"meta"`
		Data []Destination `json:"data"`
	}

	Destination struct {
		ID              int64  `json:"id"`
		Label           string `json:"label"`
		ProvinceName    string `json:"province_name"`
		CityName        string `json:"city_name"`
		DistrictName    string `json:"district_name"`
		SubdistrictName string `json:"subdistrict_name"`
		ZipCode         string `json:"zip_code"`
	}

	// ==========================================
	// Calculate Cost
	// ==========================================

	CalculateCostRequest struct {
		Origin      int64
		Destination int64
		Weight      int
		Courier     string
	}

	CalculateCostResponse struct {
		Meta Meta             `json:"meta"`
		Data []ShippingOption `json:"data"`
	}

	// Raw response standar RajaOngkir: {"rajaongkir":{"results":[...]}}
	RajaOngkirCostResponse struct {
		RajaOngkir struct {
			Results []RawCostGroup `json:"results"`
		} `json:"rajaongkir"`
	}

	RawCostGroup struct {
		Code  string          `json:"code"`
		Name  string          `json:"name"`
		Costs []RawCostService `json:"costs"`
	}

	RawCostService struct {
		Service     string        `json:"service"`
		Description string        `json:"description"`
		Cost        []RawCostItem `json:"cost"`
	}

	RawCostItem struct {
		Value int64  `json:"value"`
		Etd   string `json:"etd"`
		Note  string `json:"note"`
	}

	ShippingOption struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost        int64  `json:"cost"`
		Etd         string `json:"etd"`
	}
)
