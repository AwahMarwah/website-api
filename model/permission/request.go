package permission

type (
	CreatePermissionRequest struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Description string `json:"description"`
	}
	UpdatePermissionRequest struct {
		DisplayName string `json:"display_name" binding:"required"`
		Description string `json:"description"`
	}
)
