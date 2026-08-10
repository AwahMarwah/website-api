package cart

type (
	CartRequestBody struct {
		UserID           string `json:"user_id"`
		ProductVariantID string `binding:"required" json:"product_variant_id"`
		Qty              int    `binding:"required" json:"qty"`
	}
)
