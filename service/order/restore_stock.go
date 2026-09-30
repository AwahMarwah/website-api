package order

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"website-api/common"
	"website-api/model/order"
)

// actorInfo mencatat siapa yang mengubah status order, dipakai untuk audit trail.
type actorInfo struct {
	ID   *string
	Role *string
}

func systemActor() actorInfo {
	return actorInfo{}
}

// cancelAndRestoreStock membatalkan order PENDING sekaligus mengembalikan stok itemnya
// dalam satu transaksi.
//
// Urutan di dalamnya penting: UpdateStatusFrom dijalankan lebih dulu dengan syarat
// status = PENDING, dan baru setelah RowsAffected == 1 stok dikembalikan. Kalau dibalik,
// dua request pembatalan yang datang bersamaan bisa sama-sama mengembalikan stok.
func (s *service) cancelAndRestoreStock(orderID, targetStatus string, actor actorInfo, note string) (int, error) {
	var statusCode int

	err := s.txManager.Execute(func(tx *gorm.DB) error {
		txOrderRepo := s.orderRepo.WithTx(tx)

		affected, err := txOrderRepo.UpdateStatusFrom(orderID, common.OrderStatusPending, targetStatus)
		if err != nil {
			return fmt.Errorf("gagal memperbarui status order: %w", err)
		}
		if affected == 0 {
			statusCode = 409
			return fmt.Errorf("order sudah tidak berstatus %s", common.OrderStatusPending)
		}

		if _, err := txOrderRepo.RestoreStockForOrder(orderID); err != nil {
			return fmt.Errorf("gagal mengembalikan stok: %w", err)
		}

		if err := txOrderRepo.UpdateStatusTimestamps(orderID, targetStatus); err != nil {
			return fmt.Errorf("gagal memperbarui timestamp status: %w", err)
		}

		fromStatus := common.OrderStatusPending
		return txOrderRepo.CreateStatusHistory(order.OrderStatusHistory{
			OrderID:    orderID,
			FromStatus: &fromStatus,
			ToStatus:   targetStatus,
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
	return 200, nil
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
