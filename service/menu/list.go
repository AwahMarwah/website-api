package menu

import (
	"fmt"
	"net/http"
	"website-api/model/menu"
)

func (s *service) List(page, limit, offset int) ([]menu.MenuTree, int64, int, error) {
	all, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil menu: %w", err)
	}
	tree := buildTree(all)
	return tree, int64(len(all)), http.StatusOK, nil
}

func (s *service) Tree() ([]menu.MenuTree, int, error) {
	all, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil menu: %w", err)
	}
	return buildTree(all), http.StatusOK, nil
}
