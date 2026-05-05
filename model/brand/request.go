package brand

type (
	BrandReqQuery struct {
		Page   int `form:"page"`
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}

	FilterBrandReq struct {
		Slug string `uri:"slug" binding:"required"`
	}
)
