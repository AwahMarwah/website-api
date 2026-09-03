package menu

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func (s *service) AssignMenus(roleID string, menuIDs []string) (int, error) {
	err := s.txManager.Execute(func(tx *gorm.DB) error {
		return s.menuRepo.WithTx(tx).AssignMenus(roleID, menuIDs, time.Now())
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menetapkan menu ke role: %w", err)
	}
	return http.StatusOK, nil
}
