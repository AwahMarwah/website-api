package content_page

type (
	DetailResponse struct {
		ID      string `json:"id"`
		Slug    string `json:"slug"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  bool   `json:"status"`
	}

	FaqListResponse struct {
		Id       string `json:"id"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
		OrderNo  int    `json:"order_no"`
		IsActive bool   `json:"is_active"`
	}
)
