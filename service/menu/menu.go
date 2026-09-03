package menu

import (
	"website-api/database/transaction"
	menuModel "website-api/model/menu"
	"website-api/repository/menu"
	permissionRepo "website-api/repository/permission"
	roleRepo "website-api/repository/role"
)

type (
	IService interface {
		Create(req *menuModel.MenuCreateReq) (menuModel.MenuResponse, int, error)
		Update(id string, req *menuModel.MenuUpdateReq) (menuModel.MenuResponse, int, error)
		Delete(id string) (int, error)
		List(page, limit, offset int) ([]menuModel.MenuTree, int64, int, error)
		Tree() ([]menuModel.MenuTree, int, error)
		AssignMenus(roleID string, menuIDs []string) (int, error)
		GetMenusByRole(roleID string) ([]menuModel.MenuResponse, int, error)
		GetMyMenus(roleName string) ([]menuModel.MenuTree, int, error)
		HasPermission(roleName, permissionName string) (bool, error)
	}

	service struct {
		menuRepo       menu.IRepo
		roleRepo       roleRepo.IRepo
		permissionRepo permissionRepo.IRepo
		txManager      transaction.ITransactionManager
	}
)

func NewService(
	menuRepo menu.IRepo,
	roleRepo roleRepo.IRepo,
	permissionRepo permissionRepo.IRepo,
	txManager transaction.ITransactionManager,
) IService {
	return &service{
		menuRepo:       menuRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		txManager:      txManager,
	}
}
