package order

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// allowedTransitions memetakan status saat ini ke status yang boleh dituju (maju).
var allowedTransitions = map[string][]string{
	"PENDING":    {"PAID", "CANCELLED"},
	"PAID":       {"PROCESSING", "CANCELLED"},
	"PROCESSING": {"SHIPPED", "CANCELLED"},
	"SHIPPED":    {"COMPLETED", "CANCELLED"},
}

// UpdateStatusAdmin mengubah status order oleh admin dengan validasi transisi maju.
func (s *service) UpdateStatusAdmin(id, status string) (int, error) {
	existing, err := s.orderRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("order not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil order: %w", err)
	}

	if existing.Status == status {
		return http.StatusConflict, fmt.Errorf("status order sudah %s", status)
	}

	// status terminal tidak bisa diubah
	if existing.Status == "COMPLETED" || existing.Status == "CANCELLED" || existing.Status == "EXPIRED" {
		return http.StatusConflict, fmt.Errorf("order berstatus %s tidak dapat diubah", existing.Status)
	}

	allowed, ok := allowedTransitions[existing.Status]
	if !ok {
		return http.StatusConflict, fmt.Errorf("transisi dari status %s tidak diizinkan", existing.Status)
	}
	validTarget := false
	for _, t := range allowed {
		if t == status {
			validTarget = true
			break
		}
	}
	if !validTarget {
		return http.StatusConflict, fmt.Errorf("tidak dapat mengubah status %s menjadi %s", existing.Status, status)
	}

	if err := s.orderRepo.UpdateStatus(id, status); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui status order: %w", err)
	}
	return http.StatusOK, nil
}