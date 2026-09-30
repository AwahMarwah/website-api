package order

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"website-api/common"
)

// CancelOrder membatalkan order PENDING oleh user sendiri.
func (s *service) CancelOrder(orderID, userID, roleName string) (int, error) {
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

	if existing.Status != common.OrderStatusPending {
		return http.StatusConflict, fmt.Errorf("hanya order %s yang dapat dibatalkan", common.OrderStatusPending)
	}

	actor := actorInfo{ID: &userID, Role: &roleName}
	return s.cancelAndRestoreStock(orderID, common.OrderStatusCancelled, actor, "dibatalkan oleh pembeli")
}
