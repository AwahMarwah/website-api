package shipping

type (
	MerchantShipping struct {
		MerchantID   string          `json:"merchant_id"`
		MerchantName string          `json:"merchant_name"`
		WeightGram   int             `json:"weight_gram"`
		WeightKg     int             `json:"weight_kg"`
		Options      []ShippingOption `json:"options"`
	}

	ShippingOption struct {
		Courier string `json:"courier"`
		Service string `json:"service"`
		Cost    int64  `json:"cost"`
		Etd     string `json:"etd"`
	}
)