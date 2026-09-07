package product

type (
	ListProductResponse struct {
		Id           string  `json:"id"`
		Name         string  `json:"name"`
		Slug         string  `json:"slug"`
		BrandId      string  `json:"brand_id"`
		BrandName    string  `json:"brand_name"`
		MerchantId   string  `json:"merchant_id"`
		MerchantName string  `json:"merchant_name"`
		Thumbnail    string  `json:"thumbnail"`
		MinPrice     float64 `json:"min_price"`
		MaxPrice     float64 `json:"max_price"`
		IsInStock    int32   `json:"is_in_stock"`
		Rating       float64 `json:"rating"`
		TotalReview  int32   `json:"total_review"`
	}

	ProductListCache struct {
		Data  []ListProductResponse `json:"data"`
		Count int64                 `json:"count"`
	}

	VariantResponse struct {
		ID          string  `json:"id"`
		Sku         string  `json:"sku"`
		VariantName string  `json:"variant_name"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Weight      float32 `json:"weight"`
	}

	ProductDetailResponse struct {
		Id           string            `json:"id"`
		Name         string            `json:"name"`
		Slug         string            `json:"slug"`
		Description  string            `json:"description"`
		BrandId      string            `json:"brand_id"`
		BrandName    string            `json:"brand_name"`
		MerchantId   string            `json:"merchant_id"`
		MerchantName string            `json:"merchant_name"`
		Thumbnail    string            `json:"thumbnail"`
		BasePrice    float64           `json:"base_price"`
		MinPrice     float64           `json:"min_price"`
		MaxPrice     float64           `json:"max_price"`
		IsInStock    int32             `json:"is_in_stock"`
		Rating       float64           `json:"rating"`
		TotalReview  int32             `json:"total_review"`
		Variants     []VariantResponse `gorm:"-" json:"variants"`
	}
)
