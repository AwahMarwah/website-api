package role

import (
	"gorm.io/gorm"
	roleRepo "website-api/repository/role"
	"website-api/service/role"
)

type controller struct {
	roleService role.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{roleService: role.NewService(roleRepo.NewRepo(db))}
}
