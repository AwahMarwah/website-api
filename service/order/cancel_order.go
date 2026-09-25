package order

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// CancelOrder membatalkan order PENDING oleh user sendiri
func (s *service) CancelOrder(orderID, userID string) (int, error) {
	existing, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("order not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil order: %w", err)
	}

	// Ownership check
	if existing.UserID != userID {
		return http.StatusForbidden, fmt.Errorf("forbidden")
	}

	if existing.Status != "PENDING" {
		return http.StatusConflict, fmt.Errorf("hanya order PENDING yang dapat dibatalkan")
	}

	if err := s.orderRepo.UpdateStatus(orderID, "CANCELLED"); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membatalkan order: %w", err)
	}
	return http.StatusOK, nil
}