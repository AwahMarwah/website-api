package menu

import (
	"fmt"
	"net/http"
	"website-api/model/menu"

	"github.com/google/uuid"
)

func (s *service) Create(req *menu.MenuCreateReq) (menu.MenuResponse, int, error) {
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}

	m := &menu.Menu{
		ID:          uuid.NewString(),
		ParentID:    req.ParentID,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Icon:        req.Icon,
		Path:        req.Path,
		SortOrder:   req.SortOrder,
		IsActive:    active,
	}

	if err := s.menuRepo.Create(m); err != nil {
		return menu.MenuResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal membuat menu: %w", err)
	}
	return toResponse(*m), http.StatusCreated, nil
}
