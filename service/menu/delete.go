package menu

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func (s *service) Delete(id string) (int, error) {
	if _, err := s.menuRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("menu not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil menu: %w", err)
	}

	if err := s.menuRepo.SoftDelete(id); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus menu: %w", err)
	}
	return http.StatusOK, nil
}
