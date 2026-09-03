package menu

type (
	MenuCreateReq struct {
		ParentID    *string `json:"parent_id"`
		Name        string  `binding:"required" json:"name"`
		DisplayName string  `binding:"required" json:"display_name"`
		Icon        string  `json:"icon"`
		Path        string  `json:"path"`
		SortOrder   int     `json:"sort_order"`
		IsActive    *bool   `json:"is_active"`
	}

	MenuUpdateReq struct {
		ParentID    *string `json:"parent_id"`
		Name        string  `json:"name"`
		DisplayName string  `json:"display_name"`
		Icon        string  `json:"icon"`
		Path        string  `json:"path"`
		SortOrder   int     `json:"sort_order"`
		IsActive    *bool   `json:"is_active"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	AssignMenusReq struct {
		MenuIDs []string `binding:"required" json:"menu_ids"`
	}

	ListMenuReqQuery struct {
		Page   int `form:"page"`
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}
)
