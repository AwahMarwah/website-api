package menu

import (
	"errors"
	"fmt"
	"net/http"
	roleModel "website-api/model/role"
	"website-api/model/menu"

	"gorm.io/gorm"
)

func (s *service) GetMenusByRole(roleID string) ([]menu.MenuResponse, int, error) {
	if _, err := s.roleRepo.Take([]string{"id"}, &roleModel.Role{Id: roleID}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, fmt.Errorf("role not found")
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil role: %w", err)
	}

	menus, err := s.menuRepo.FindMenusByRole(roleID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil menu role: %w", err)
	}

	res := make([]menu.MenuResponse, 0)
	for _, m := range menus {
		res = append(res, toResponse(m))
	}
	return res, http.StatusOK, nil
}
