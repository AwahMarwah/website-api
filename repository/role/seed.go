package role

import "website-api/model/role"

func (r *repo) Seed() error {
	return r.db.Create([]role.Role{
		{
			Name:        "super_admin",
			DisplayName: "Super Administrator",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
		{
			Name:        "staff",
			DisplayName: "Staff",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
		{
			Name:        "finance",
			DisplayName: "Finance",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
		{
			Name:        "merchant",
			DisplayName: "Merchant",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
		{
			Name:        "consumer",
			DisplayName: "Consumer",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
		{
			Name:        "guest",
			DisplayName: "Guest",
			Description: "",
			IsActive:    true,
			CreatedBy:   "system",
		},
	}).Error
}
