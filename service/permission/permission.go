package permission

import (
	"fmt"
	permissionModel "website-api/model/permission"
	menuModel "website-api/model/menu"
	permissionRepo "website-api/repository/permission"
)

type (
	IService interface {
		List() ([]menuModel.PermissionResponse, error)
	}

	service struct {
		permissionRepo permissionRepo.IRepo
	}
)

func NewService(permissionRepo permissionRepo.IRepo) IService {
	return &service{permissionRepo: permissionRepo}
}

func (s *service) List() ([]menuModel.PermissionResponse, error) {
	perms, err := s.permissionRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil permission: %w", err)
	}
	res := make([]menuModel.PermissionResponse, 0)
	for _, p := range perms {
		res = append(res, mapToResponse(p))
	}
	return res, nil
}

func mapToResponse(p permissionModel.Permission) menuModel.PermissionResponse {
	return menuModel.PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		MenuID:      p.MenuId,
	}
}
