package permission

import (
	permissionRepo "website-api/repository/permission"
	"website-api/service/permission"

	"gorm.io/gorm"
)

type controller struct {
	permissionService permission.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{permissionService: permission.NewService(permissionRepo.NewRepo(db))}
}
