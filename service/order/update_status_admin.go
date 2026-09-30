package order

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"website-api/common"
	"website-api/model/order"
)

// allowedTransitions memetakan status saat ini ke status yang boleh dituju (maju).
var allowedTransitions = map[string][]string{
	common.OrderStatusPending:    {common.OrderStatusPaid, common.OrderStatusCancelled, common.OrderStatusExpired},
	common.OrderStatusPaid:       {common.OrderStatusProcessing, common.OrderStatusCancelled},
	common.OrderStatusProcessing: {common.OrderStatusShipped, common.OrderStatusCancelled},
	common.OrderStatusShipped:    {common.OrderStatusCompleted, common.OrderStatusCancelled},
}

// UpdateStatusAdmin mengubah status order oleh admin dengan validasi transisi maju.
// Pembatalan dari status non-PENDING ikut mengembalikan stok karena inventory sudah
// dikurangi sejak checkout; hanya konfirmasi pembayaran yang tidak mengembalikannya.
func (s *service) UpdateStatusAdmin(id, actorID, actorRole, status, note string) (int, error) {
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
	if existing.Status == common.OrderStatusCompleted ||
		existing.Status == common.OrderStatusCancelled ||
		existing.Status == common.OrderStatusExpired {
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

	actor := actorInfo{ID: &actorID, Role: &actorRole}

	// Batal setelah pembayaran: stok harus kembali. Jalankan lewat helper yang sama
	// dengan pembatalan oleh pembeli supaya perilakunya identik.
	if status == common.OrderStatusCancelled {
		if existing.Status == common.OrderStatusPending {
			return s.cancelAndRestoreStock(id, common.OrderStatusCancelled, actor, note)
		}
		return s.cancelPaidAndRestoreStock(id, existing.Status, actor, note)
	}

	if err := s.applyStatusChange(id, existing.Status, status, actor, note); err != nil {
		return http.StatusInternalServerError, err
	}

	// Komisi marketplace dicatat saat order selesai, bukan saat dibayar, supaya
	// refund/retur tidak butuh koreksi komisi manual.
	if status == common.OrderStatusCompleted {
		s.recordCommissionAfterCompletion(id)
	}
	return http.StatusOK, nil
}

// cancelPaidAndRestoreStock membatalkan order yang sudah dibayar (mis. returSeller,
// permintaan pembeli) dan mengembalikan stok.
func (s *service) cancelPaidAndRestoreStock(id, fromStatus string, actor actorInfo, note string) (int, error) {
	var statusCode int

	err := s.txManager.Execute(func(tx *gorm.DB) error {
		txOrderRepo := s.orderRepo.WithTx(tx)

		affected, err := txOrderRepo.UpdateStatusFrom(id, fromStatus, common.OrderStatusCancelled)
		if err != nil {
			return fmt.Errorf("gagal memperbarui status order: %w", err)
		}
		if affected == 0 {
			statusCode = 409
			return fmt.Errorf("status order sudah berubah")
		}

		if _, err := txOrderRepo.RestoreStockForOrder(id); err != nil {
			return fmt.Errorf("gagal mengembalikan stok: %w", err)
		}
		if err := txOrderRepo.UpdateStatusTimestamps(id, common.OrderStatusCancelled); err != nil {
			return fmt.Errorf("gagal memperbarui timestamp status: %w", err)
		}

		from := fromStatus
		return txOrderRepo.CreateStatusHistory(order.OrderStatusHistory{
			OrderID:    id,
			FromStatus: &from,
			ToStatus:   common.OrderStatusCancelled,
			ActorID:    actor.ID,
			ActorRole:  actor.Role,
			Note:       optionalString(note),
			CreatedAt:  time.Now(),
		})
	})
	if err != nil {
		if statusCode == 0 {
			statusCode = 500
		}
		return statusCode, err
	}
	return http.StatusOK, nil
}

// applyStatusChange menyimpan perubahan status beserta timestamp dan audit trail.
func (s *service) applyStatusChange(orderID, fromStatus, toStatus string, actor actorInfo, note string) error {
	return s.txManager.Execute(func(tx *gorm.DB) error {
		txOrderRepo := s.orderRepo.WithTx(tx)

		affected, err := txOrderRepo.UpdateStatusFrom(orderID, fromStatus, toStatus)
		if err != nil {
			return fmt.Errorf("gagal memperbarui status order: %w", err)
		}
		if affected == 0 {
			return fmt.Errorf("status order sudah berubah oleh request lain")
		}
		if err := txOrderRepo.UpdateStatusTimestamps(orderID, toStatus); err != nil {
			return fmt.Errorf("gagal memperbarui timestamp status: %w", err)
		}

		from := fromStatus
		return txOrderRepo.CreateStatusHistory(order.OrderStatusHistory{
			OrderID:    orderID,
			FromStatus: &from,
			ToStatus:   toStatus,
			ActorID:    actor.ID,
			ActorRole:  actor.Role,
			Note:       optionalString(note),
			CreatedAt:  time.Now(),
		})
	})
}
