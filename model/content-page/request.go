package content_page

type (
	ReqPath struct {
		Slug string `binding:"required" uri:"slug"`
	}

	FaqListReqQuery struct {
		Limit  int `form:"limit"`
		Offset int
		Page   int `form:"page"`
	}
)
