package menu

import (
	"errors"
	"fmt"
	"net/http"
	"website-api/model/menu"

	"gorm.io/gorm"
)

func (s *service) Update(id string, req *menu.MenuUpdateReq) (menu.MenuResponse, int, error) {
	existing, err := s.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menu.MenuResponse{}, http.StatusNotFound, fmt.Errorf("menu not found")
		}
		return menu.MenuResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengambil menu: %w", err)
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.DisplayName != "" {
		existing.DisplayName = req.DisplayName
	}
	if req.Icon != "" {
		existing.Icon = req.Icon
	}
	if req.Path != "" {
		existing.Path = req.Path
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}
	existing.SortOrder = req.SortOrder
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.menuRepo.Update(&existing); err != nil {
		return menu.MenuResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal memperbarui menu: %w", err)
	}
	return toResponse(existing), http.StatusOK, nil
}
