package role

type (
	RoleCreateReqBody struct {
		Name        string `binding:"required" json:"name"`
		DisplayName string `binding:"required" json:"display_name"`
		Description string `binding:"required" json:"description"`
		IsActive    bool   `binding:"required" json:"is_active"`
	}

	ReqPath struct {
		Id string `uri:"id" binding:"required"`
	}

	UpdateReq struct {
		Path struct {
			Id string `uri:"id" binding:"required"`
		}
		Body struct {
			Name        string `binding:"required" json:"name"`
			DisplayName string `binding:"required" json:"display_name"`
			Description string `binding:"required" json:"description"`
			IsActive    *bool  `binding:"required" json:"is_active"`
		}
	}
)
