package order

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	settlementService "website-api/service/settlement"
)

// recordCommission menghitung dan mencatat komisi marketplace untuk setiap merchant
// di dalam order yang sudah selesai.
//
// Komisi dicatat saat COMPLETED, bukan saat PAID. Kalau dicatat saat PAID, setiap
// refund atau retur jadi butuh koreksi manual. Saat COMPLETED, pembeli sudah menerima
// barang sehingga komisi dianggap sudahearned.
//
// Kegagalan di sini tidak menggagalkan perubahan status order: RecordCommission
// bersifat idempoten, jadi retry atau cron rekonsiliasi bisa memanggilnya lagi
// tanpa risiko komisi terhitung dua kali.
func (s *service) recordCommission(orderID string, tx *gorm.DB) error {
	o, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("gagal mengambil order: %w", err)
	}

	// Kelompokkan item order per merchant: satu order bisa berisi barang banyak seller.
	itemsByMerchant := make(map[string][]settlementService.CommissionItem)
	for _, item := range o.Items {
		if item.MerchantID == "" {
			// Item tanpa merchant tidak punya komisi; dilewati agar tidak menggagalkan
			// pencatatan komisi seller lain.
			continue
		}
		itemsByMerchant[item.MerchantID] = append(itemsByMerchant[item.MerchantID], settlementService.CommissionItem{
			OrderItemID: item.ID,
			Subtotal:    item.Subtotal,
		})
	}

	for merchantID, items := range itemsByMerchant {
		m, err := s.merchantRepo.FindByID(merchantID)
		if err != nil {
			return fmt.Errorf("gagal mengambil merchant %s: %w", merchantID, err)
		}
		if err := s.settlementService.RecordCommission(tx, merchantID, orderID, m.CommissionRateBP, items); err != nil {
			return fmt.Errorf("gagal mencatat komisi merchant %s: %w", merchantID, err)
		}
	}
	return nil
}

// recordCommissionAfterCompletion mencatat komisi dalam transaksi terpisah di luar
// perubahan status order, supaya kegagalan komisi tidak membatalkan status order yang
// sebenarnya sudah berhasil diubah.
func (s *service) recordCommissionAfterCompletion(orderID string) {
	if s.settlementService == nil {
		return
	}
	if err := s.txManager.Execute(func(tx *gorm.DB) error {
		return s.recordCommission(orderID, tx)
	}); err != nil {
		log.Printf("gagal mencatat komisi untuk order %s: %v", orderID, err)
	}
}
