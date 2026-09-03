package menu

type (
	MenuResponse struct {
		ID          string  `json:"id"`
		ParentID    *string `json:"parent_id"`
		Name        string  `json:"name"`
		DisplayName string  `json:"display_name"`
		Icon        string  `json:"icon"`
		Path        string  `json:"path"`
		SortOrder   int     `json:"sort_order"`
		IsActive    bool    `json:"is_active"`
	}

	MenuTree struct {
		ID          string      `json:"id"`
		ParentID    *string     `json:"parent_id"`
		Name        string      `json:"name"`
		DisplayName string      `json:"display_name"`
		Icon        string      `json:"icon"`
		Path        string      `json:"path"`
		SortOrder   int         `json:"sort_order"`
		IsActive    bool        `json:"is_active"`
		Children    []MenuTree  `json:"children"`
	}

	PermissionResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
		MenuID      *string `json:"menu_id"`
	}
)
