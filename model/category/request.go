package category

type (
	FilterCategory struct {
		Search string `form:"search"`
		Page   int    `form:"page"`
		Limit  int    `form:"limit"`
		Offset int    `form:"offset"`
	}
)
