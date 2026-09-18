package product

type (
	GetListProductReqQuerry struct {
		Page     int     `form:"page"`
		Limit    int     `form:"limit"`
		Offset   int     `form:"offset"`
		Brand    string  `form:"brand"`
		Category string  `form:"category"`
		MinPrice float64 `form:"min_price"`
		MaxPrice float64 `form:"max_price"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	CreateProductReq struct {
		Name        string   `binding:"required" json:"name"`
		BrandId     string   `binding:"required" json:"brand_id"`
		MerchantId  string   `json:"merchant_id"`
		Sku         string   `binding:"required" json:"sku"`
		Slug        string   `binding:"required" json:"slug"`
		Description string   `json:"description"`
		BasePrice   float64  `json:"base_price"`
		Categories  []string `json:"categories"`
	}

	UpdateProductReq struct {
		Name        string   `json:"name"`
		BrandId     string   `json:"brand_id"`
		MerchantId  string   `json:"merchant_id"`
		Sku         string   `json:"sku"`
		Slug        string   `json:"slug"`
		Description string   `json:"description"`
		BasePrice   float64  `json:"base_price"`
		Status      string   `json:"status"`
		Categories  []string `json:"categories"`
	}
)
