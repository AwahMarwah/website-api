package merchant

import "time"

type (
	MerchantResponse struct {
		ID               string  `json:"id"`
		Name             string  `json:"name"`
		Slug             string  `json:"slug"`
		DestinationID    int64   `json:"destination_id"`
		CityID           string  `json:"city_id"`
		Address          string  `json:"address"`
		IsActive         bool    `json:"is_active"`
		UserID           string  `json:"user_id"`
		CommissionRateBP int     `json:"commission_rate_bp"`
	}

	SellerVariantResponse struct {
		ID          string  `json:"id"`
		Sku         string  `json:"sku"`
		VariantName string  `json:"variant_name"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Weight      float32 `json:"weight"`
		IsActive    bool    `json:"is_active"`
	}

	SellerProductResponse struct {
		ID          string                  `json:"id"`
		Name        string                  `json:"name"`
		Slug        string                  `json:"slug"`
		Description string                  `json:"description"`
		BrandId     string                  `json:"brand_id"`
		Status      string                  `json:"status"`
		BasePrice   float32                 `json:"base_price"`
		WeightGram  int                     `json:"weight_gram"`
		Thumbnail   string                  `json:"thumbnail"`
		CreatedAt   time.Time               `json:"created_at"`
		UpdatedAt   time.Time               `json:"updated_at"`
		Variants    []SellerVariantResponse `json:"variants"`
	}

	// SellerOrderItemResponse hanya memuat item milik seller yang meminta.
	// Seller tidak boleh melihat item seller lain di order yang sama.
	SellerOrderItemResponse struct {
		ID              string  `json:"id"`
		ProductID       string  `json:"product_id"`
		ProductName     string  `json:"product_name"`
		VariantName     string  `json:"variant_name"`
		ProductImageURL string  `json:"product_image_url"`
		Sku             string  `json:"sku"`
		ProductVariantID string `json:"product_variant_id"`
		Price           float64 `json:"price"`
		Qty             int     `json:"qty"`
		Subtotal        float64 `json:"subtotal"`
		TotalWeightGram int     `json:"total_weight_gram"`
	}

	SellerOrderResponse struct {
		ID              string                        `json:"id"`
		BuyerName       string                        `json:"buyer_name"`
		Status          string                        `json:"status"`
		ItemSubtotal    float64                       `json:"item_subtotal"`
		ShippingCost    int64                         `json:"shipping_cost"`
		TotalAmount     float64                       `json:"total_amount"`
		Courier         string                        `json:"courier"`
		Service         string                        `json:"service"`
		Etd             string                        `json:"etd"`
		CreatedAt       time.Time                     `json:"created_at"`
		PaidAt          *time.Time                    `json:"paid_at"`
		ShippedAt       *time.Time                    `json:"shipped_at"`
		CompletedAt     *time.Time                    `json:"completed_at"`
		Items           []SellerOrderItemResponse     `json:"items"`
	}

	SellerStatsResponse struct {
		TotalProducts    int64   `json:"total_products"`
		ActiveProducts   int64   `json:"active_products"`
		DraftProducts    int64   `json:"draft_products"`
		TotalOrders      int64   `json:"total_orders"`
		PendingOrders    int64   `json:"pending_orders"`
		ProcessingOrders int64   `json:"processing_orders"`
		ShippedOrders    int64   `json:"shipped_orders"`
		CompletedOrders  int64   `json:"completed_orders"`
		ItemRevenue      float64 `json:"item_revenue"`
		CommissionEarned float64 `json:"commission_earned"`
		Balance          float64 `json:"balance"`
	}
)
