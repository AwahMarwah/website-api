package brand

type (
	ListBrandResponse struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		LogoUrl   string `json:"logo_url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)
