package voucher

type (
	CreateVoucherReq struct {
		Code         string   `binding:"required" json:"code"`
		Description  string   `json:"description"`
		Type         string   `binding:"required,oneof=PERCENT FIXED" json:"type"`
		Value        float64  `binding:"required,gt=0" json:"value"`
		MaxDiscount  *float64 `json:"max_discount"`
		MinSpend     float64  `json:"min_spend"`
		Quota        *int     `json:"quota"`
		PerUserLimit int      `json:"per_user_limit"`
		StartsAt     *string  `json:"starts_at"`
		EndsAt       *string  `json:"ends_at"`
		MerchantID   *string  `json:"merchant_id"`
		ProductIDs   []string `json:"product_ids"`
	}

	UpdateVoucherReq struct {
		Description  *string  `json:"description"`
		Value        *float64 `json:"value"`
		MaxDiscount  *float64 `json:"max_discount"`
		MinSpend     *float64 `json:"min_spend"`
		Quota        *int     `json:"quota"`
		PerUserLimit *int     `json:"per_user_limit"`
		StartsAt     *string  `json:"starts_at"`
		EndsAt       *string  `json:"ends_at"`
		IsActive     *bool    `json:"is_active"`
		ProductIDs   []string `json:"product_ids"`
	}

	ValidateVoucherReqQuery struct {
		Code   string  `form:"code" binding:"required"`
		Amount float64 `form:"amount"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	ListVoucherReqQuery struct {
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Offset     int    `form:"offset"`
		Search     string `form:"search"`
		IsActive   *bool  `form:"is_active"`
		MerchantID string `form:"merchant_id"`
	}
)
