package role

import "time"

type (
	DetailResponse struct {
		Id          string     `json:"id"`
		Name        string     `json:"name"`
		DisplayName string     `json:"display_name"`
		Description string     `json:"description"`
		IsActive    bool       `json:"is_active"`
		CreatedAt   time.Time  `json:"created_at"`
		CreatedBy   string     `json:"created_by"`
		UpdatedAt   *time.Time `json:"updated_at"`
		UpdatedBy   string     `json:"updated_by"`
		DeletedAt   *time.Time `json:"deleted_at"`
		DeletedBy   string     `json:"deleted_by"`
	}
)
