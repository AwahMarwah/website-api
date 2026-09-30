package order

type (
	CheckoutRequest struct {
		UserID        string         `json:"-"`
		AddressID     string         `binding:"required" json:"address_id"`
		PaymentMethod string         `binding:"required" json:"payment_method"`
		ShippingFee   float64        `binding:"required" json:"shipping_fee"`
		Items         []CheckoutItem `binding:"required,min=1" json:"items"`
		Shippings     []ShippingsReq `json:"shippings"`
		VoucherCode   string         `json:"voucher_code"`
		Note          string         `json:"note"`
	}

	ShippingsReq struct {
		MerchantID string `json:"merchant_id"`
		Courier    string `json:"courier"`
		Service    string `json:"service"`
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

	UpdateOrderStatusReq struct {
		Status string `binding:"required,oneof=PENDING PAID PROCESSING SHIPPED COMPLETED CANCELLED EXPIRED" json:"status"`
		Note   string `json:"note"`
	}
)
