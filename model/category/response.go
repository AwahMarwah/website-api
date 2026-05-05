package category

type (
	ListCategoryResponse struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		ParentId  string `json:"parent_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)
