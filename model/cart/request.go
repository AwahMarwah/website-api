package cart

type (
	CartRequestBody struct {
		UserID           string `json:"user_id"`
		ProductVariantID string `binding:"required" json:"product_variant_id"`
		Qty              int    `binding:"required" json:"qty"`
	}

	DeleteCartRequest struct {
		UserID     string   `json:"user_id"`
		VariantIDs []string `binding:"required,gt=0"json:"variant_ids"`
	}
)
