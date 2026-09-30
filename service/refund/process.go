package refund

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"website-api/common"
	"website-api/model/order"
	refundModel "website-api/model/refund"
)

// Approve menandai permintaan refund disetujui finance. Statusnya masih APPROVED:
// dana baru benar-benar dikembalikan saat Process dipanggil.
func (s *service) Approve(refundID, approverID string) (int, error) {
	existing, err := s.refundRepo.FindByID(refundID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("refund not found")
	}
	if existing.Status != refundModel.StatusPending {
		return http.StatusConflict, fmt.Errorf("refund sudah diproses dengan status %s", existing.Status)
	}

	if err := s.refundRepo.UpdateStatus(refundID, refundModel.StatusApproved, map[string]any{
		"approved_by": &approverID,
		"updated_at":  time.Now(),
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menyetujui refund: %w", err)
	}
	return http.StatusOK, nil
}

// Reject menolak permintaan refund. Tidak ada efek samping karena dana belum berpindah.
func (s *service) Reject(refundID, approverID, reason string) (int, error) {
	existing, err := s.refundRepo.FindByID(refundID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("refund not found")
	}
	if existing.Status == refundModel.StatusCompleted {
		return http.StatusConflict, fmt.Errorf("refund yang sudah selesai tidak bisa ditolak")
	}

	values := map[string]any{
		"approved_by": &approverID,
		"updated_at":  time.Now(),
	}
	if reason != "" {
		values["failure_reason"] = &reason
	}
	if err := s.refundRepo.UpdateStatus(refundID, refundModel.StatusRejected, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menolak refund: %w", err)
	}
	return http.StatusOK, nil
}

// Process mengirim refund ke payment provider lalu mengembalikan stok dan mengoreksi
// saldo seller.
//
// Urutan di dalamnya penting: panggilan ke provider dilakukan DI LUAR transaksi.
// Panggilan HTTP tidak bisa di-rollback, jadi kalau refund berhasil di provider tapi
// commit database gagal, dana sudah keluar sementara stok belum kembali. Karena itu
// statusnya diubah ke PROCESSING lebih dulu, supaya refund yang gagal di tengah
// masih terlihat dan bisa diproses ulang.
func (s *service) Process(refundID string) (int, error) {
	existing, err := s.refundRepo.FindByID(refundID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("refund not found")
	}
	if existing.Status != refundModel.StatusApproved {
		return http.StatusConflict, fmt.Errorf("refund harus disetujui sebelum diproses")
	}

	// Tandai sedang diproses supaya tidak ada dua worker yang memanggil provider
	// untuk refund yang sama.
	affected, err := s.claimForProcessing(refundID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengunci refund: %w", err)
	}
	if affected == 0 {
		return http.StatusConflict, fmt.Errorf("refund sedang diproses")
	}

	o, err := s.orderRepo.FindByID(existing.OrderID)
	if err != nil {
		s.markFailed(refundID, "order tidak ditemukan")
		return http.StatusNotFound, fmt.Errorf("order not found")
	}

	providerResp, err := s.provider.Refund(o.ID, transactionStatusFor(o.Status), reasonOr(existing.Reason), int64(existing.Amount))
	if err != nil {
		s.markFailed(refundID, err.Error())
		return http.StatusBadGateway, fmt.Errorf("gagal memproses refund ke payment provider: %w", err)
	}

	if err := s.finalize(refundID, providerResp); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *service) claimForProcessing(refundID string) (int64, error) {
	var affected int64
	err := s.txManager.Execute(func(tx *gorm.DB) error {
		result := tx.Model(&refundModel.Refund{}).
			Where("id = ? AND status = ?", refundID, refundModel.StatusApproved).
			Updates(map[string]any{
				"status":     refundModel.StatusProcessing,
				"updated_at": time.Now(),
			})
		affected = result.RowsAffected
		return result.Error
	})
	return affected, err
}

func (s *service) markFailed(refundID, reason string) {
	_ = s.refundRepo.UpdateStatus(refundID, refundModel.StatusFailed, map[string]any{
		"failure_reason": &reason,
		"updated_at":     time.Now(),
	})
}

// finalize mencatat refund berhasil: mengembalikan stok item, mengembalikan komisi
// yang sudah dipotong ke seller, dan menandai order sebagai REFUNDED.
func (s *service) finalize(refundID string, providerResp *refundResponse) error {
	return s.txManager.Execute(func(tx *gorm.DB) error {
		txRefundRepo := s.refundRepo.WithTx(tx)
		txOrderRepo := s.orderRepo.WithTx(tx)

		refundRow, err := txRefundRepo.FindByID(refundID)
		if err != nil {
			return err
		}

		// Stok dikembalikan hanya sekali: flag restocked dicek dalam UPDATE ber-syarat.
		if !refundRow.Restocked {
			if _, err := txOrderRepo.RestoreStockForOrder(refundRow.OrderID); err != nil {
				return fmt.Errorf("gagal mengembalikan stok: %w", err)
			}
		}

		// Kembalikan komisi yang sudah dipotong dari saldo seller. Tanpa ini seller
		// tetap kehilangan komisi dari order yang uangnya sudah dikembalikan ke pembeli.
		commission, err := s.commissionReversed(tx, refundRow)
		if err != nil {
			return err
		}
		if commission > 0 && refundRow.MerchantID != nil {
			note := "koreksi komisi karena refund"
			if err := s.settlementService.RecordRefund(tx, *refundRow.MerchantID, refundRow.OrderID, note, commission); err != nil {
				return fmt.Errorf("gagal mengoreksi saldo seller: %w", err)
			}
		}

		values := map[string]any{
			"restocked":  true,
			"updated_at": time.Now(),
		}
		if providerResp != nil && providerResp.RefundID != "" {
			values["provider_refund_id"] = providerResp.RefundID
		}
		if err := txRefundRepo.UpdateStatus(refundID, refundModel.StatusCompleted, values); err != nil {
			return err
		}

		// Order ditandai REFUNDED hanya kalau seluruh nilainya sudah dikembalikan;
		// refund parsial tetap mempertahankan status order sebelumnya.
		refunded, err := txRefundRepo.SumCompletedByOrder(refundRow.OrderID)
		if err != nil {
			return err
		}
		o, err := txOrderRepo.FindByID(refundRow.OrderID)
		if err != nil {
			return err
		}
		if refunded >= o.TotalAmount {
			from := o.Status
			if _, err := txOrderRepo.UpdateStatusFrom(o.ID, o.Status, common.OrderStatusRefunded); err != nil {
				return err
			}
			if err := txOrderRepo.UpdateStatusTimestamps(o.ID, common.OrderStatusRefunded); err != nil {
				return err
			}
			return txOrderRepo.CreateStatusHistory(order.OrderStatusHistory{
				OrderID:    o.ID,
				FromStatus: &from,
				ToStatus:   common.OrderStatusRefunded,
				ActorRole:  optionalString("finance"),
				Note:       optionalString("dikembalikan ke pembeli"),
				CreatedAt:  time.Now(),
			})
		}
		return nil
	})
}

func transactionStatusFor(orderStatus string) string {
	switch orderStatus {
	case common.OrderStatusPaid, common.OrderStatusProcessing:
		return "capture"
	case common.OrderStatusShipped, common.OrderStatusCompleted:
		return "settlement"
	}
	return "settlement"
}

func reasonOr(v *string) string {
	if v == nil {
		return "pengembalian dana"
	}
	return *v
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
