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

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	CreateBrandReq struct {
		Name    string `binding:"required" json:"name"`
		Slug    string `binding:"required" json:"slug"`
		LogoUrl string `json:"logo_url"`
	}

	UpdateBrandReq struct {
		Name    string `json:"name"`
		Slug    string `json:"slug"`
		LogoUrl string `json:"logo_url"`
	}
)
