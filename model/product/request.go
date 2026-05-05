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
)
