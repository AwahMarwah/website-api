package menu

import (
	"errors"
	"fmt"
	roleModel "website-api/model/role"

	"gorm.io/gorm"
)

// HasPermission memeriksa apakah role memiliki permission tertentu.
// super_admin dan admin selalu memiliki semua permission (bypass RBAC).
func (s *service) HasPermission(roleName, permissionName string) (bool, error) {
	if roleName == "super_admin" || roleName == "admin" {
		return true, nil
	}

	role, err := s.roleRepo.Take([]string{"id"}, &roleModel.Role{Name: roleName})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("role not found")
		}
		return false, err
	}

	names, err := s.permissionRepo.FindNamesByRole(role.Id)
	if err != nil {
		return false, err
	}

	for _, n := range names {
		if n == permissionName {
			return true, nil
		}
	}
	return false, nil
}
