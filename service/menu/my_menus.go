package menu

import (
	"errors"
	"fmt"
	"net/http"
	roleModel "website-api/model/role"
	"website-api/model/menu"

	"gorm.io/gorm"
)

func (s *service) GetMyMenus(roleName string) ([]menu.MenuTree, int, error) {
	role, err := s.roleRepo.Take([]string{"id"}, &roleModel.Role{Name: roleName})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, fmt.Errorf("role not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil role: %w", err)
	}

	menus, err := s.menuRepo.FindMenusByRole(role.Id)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil menu: %w", err)
	}

	return buildTree(menus), http.StatusOK, nil
}
