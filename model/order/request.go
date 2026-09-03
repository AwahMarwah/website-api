package order

type (
	CheckoutRequest struct {
		UserID        string         `json:"user_id"`
		AddressID     string         `binding:"required" json:"address_id"`
		PaymentMethod string         `binding:"required" json:"payment_method"`
		ShippingFee   float64        `binding:"required" json:"shipping_fee"`
		Items         []CheckoutItem `json:"items"`
	}

	CheckoutItem struct {
		VariantID string `json:"variant_id"`
		Qty       int    `json:"qty"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	ListOrderReqQuery struct {
		Page     int    `form:"page"`
		Limit    int    `form:"limit"`
		Offset   int    `form:"offset"`
		Status   string `form:"status"`
		UserID   string `json:"-"`
		IsAdmin  bool   `json:"-"`
	}
)
