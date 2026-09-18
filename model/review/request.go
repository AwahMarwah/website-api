package review

type (
	CreateReviewReq struct {
		ProductID string `json:"product_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required,min=1,max=5"`
		Comment   string `json:"comment"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	ListReviewReqQuery struct {
		Page   int `form:"page"`
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}
)