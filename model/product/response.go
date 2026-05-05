package product

type (
	ListProductResponse struct {
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Slug        string  `json:"slug"`
		BrandId     string  `json:"brand_id"`
		BrandName   string  `json:"brand_name"`
		Thumbnail   string  `json:"thumbnail"`
		MinPrice    float64 `json:"min_price"`
		MaxPrice    float64 `json:"max_price"`
		IsInStock   int32   `json:"is_in_stock"`
		Rating      float64 `json:"rating"`
		TotalReview int32   `json:"total_review"`
	}

	ProductListCache struct {
		Data  []ListProductResponse `json:"data"`
		Count int64                 `json:"count"`
	}
)
