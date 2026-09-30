package merchant

type (
	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	CreateMerchantReq struct {
		Name          string `binding:"required" json:"name"`
		Slug          string `binding:"required" json:"slug"`
		DestinationID int64  `binding:"required" json:"destination_id"`
		CityID        string `json:"city_id"`
		Address       string `json:"address"`
	}

	UpdateMerchantReq struct {
		Name          string `json:"name"`
		Slug          string `json:"slug"`
		DestinationID *int64 `json:"destination_id"`
		CityID        string `json:"city_id"`
		Address       string `json:"address"`
		IsActive      *bool  `json:"is_active"`
	}

	ListSellerProductReqQuery struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
		Search string `form:"search"`
		Status string `form:"status" binding:"omitempty,oneof=draft active inactive"`
		Sort   string `form:"sort"`
	}

	ListSellerOrderReqQuery struct {
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
		Status string `form:"status" binding:"omitempty,oneof=PENDING PAID PROCESSING SHIPPED COMPLETED CANCELLED EXPIRED REFUNDED"`
	}

	// CreateSellerProductReq tidak punya field MerchantId: pemilik produk selalu
	// diturunkan dari session seller, tidak pernah dari body request.
	CreateSellerProductReq struct {
		Name        string   `binding:"required" json:"name"`
		BrandId     string   `binding:"required" json:"brand_id"`
		Sku         string   `binding:"required" json:"sku"`
		Slug        string   `binding:"required" json:"slug"`
		Description string   `json:"description"`
		BasePrice   float64  `json:"base_price"`
		WeightGram  int      `json:"weight_gram"`
		Categories  []string `json:"categories"`
	}

	UpdateSellerProductReq struct {
		Name        string   `json:"name"`
		BrandId     *string  `json:"brand_id"`
		Slug        string   `json:"slug"`
		Description string   `json:"description"`
		BasePrice   *float64 `json:"base_price"`
		WeightGram  *int     `json:"weight_gram"`
		Categories  []string `json:"categories"`
	}

	CreateSellerVariantReq struct {
		Sku         string  `binding:"required" json:"sku"`
		VariantName string  `binding:"required" json:"variant_name"`
		Price       float64 `binding:"required,gt=0" json:"price"`
		Stock       int     `binding:"gte=0" json:"stock"`
		Weight      float32 `json:"weight"`
	}

	UpdateSellerVariantReq struct {
		Sku         *string  `json:"sku"`
		VariantName *string  `json:"variant_name"`
		Price       *float64 `json:"price"`
		Stock       *int     `json:"stock"`
		Weight      *float32 `json:"weight"`
		IsActive    *bool    `json:"is_active"`
	}
)
