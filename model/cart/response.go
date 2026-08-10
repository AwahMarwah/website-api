package cart

type (
	CartItemResponse struct {
		CartItemID        string  `json:"cart_item_id"`
		ProductID         string  `json:"product_id"`
		ProductName       string  `json:"product_name"`
		ProductVariantID  string  `json:"product_variant_id"`
		VariantName       string  `json:"variant_name"`
		Price             float64 `json:"price"`
		FormattedPrice    string  `json:"formatted_price"`
		Qty               int     `json:"qty"`
		SubTotal          float64 `json:"sub_total"`
		FormattedSubTotal string  `json:"formatted_sub_total"`
		CurrentStock      int     `json:"current_stock"`
	}

	Summary struct {
		TotalItems          int     `json:"total_items"`
		GrandTotal          float64 `json:"grand_total"`
		FormattedGrandTotal string  `json:"formatted_grand_total"`
	}

	GetCartByUserIDResponse struct {
		UserID  string             `json:"user_id"`
		Items   []CartItemResponse `json:"items"`
		Summary Summary            `json:"summary"`
	}
)
