package category

type (
	FilterCategory struct {
		Search string `form:"search"`
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	CategorySlugPath struct {
		Slug string `uri:"slug" binding:"required"`
	}

	CreateCategoryReq struct {
		Name     string `binding:"required" json:"name"`
		Slug     string `binding:"required" json:"slug"`
		ParentId string `json:"parent_id"`
	}

	UpdateCategoryReq struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		ParentId string `json:"parent_id"`
	}
)
