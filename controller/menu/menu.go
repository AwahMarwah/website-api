package menu

import (
	"website-api/database/transaction"
	menuRepo "website-api/repository/menu"
	permissionRepo "website-api/repository/permission"
	roleRepo "website-api/repository/role"
	"website-api/service/menu"

	"gorm.io/gorm"
)

type controller struct {
	menuService menu.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{
		menuService: menu.NewService(
			menuRepo.NewRepo(db),
			roleRepo.NewRepo(db),
			permissionRepo.NewRepo(db),
			transaction.NewTransactionManager(db),
		),
	}
}