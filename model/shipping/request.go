package shipping

type (
	ShippingCostRequest struct {
		AddressID string             `binding:"required" json:"address_id"`
		Items     []ShippingCostItem `binding:"required,dive" json:"items"`
		Couriers  []string           `json:"couriers"`
	}

	ShippingCostItem struct {
		VariantID string `binding:"required" json:"variant_id"`
		Qty       int    `binding:"required,min=1" json:"qty"`
	}
)