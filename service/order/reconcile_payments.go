package order

import (
	"errors"
	"log"
	"strconv"

	midtransProvider "website-api/third-party/provider/midtrans"
)

var errEmptyProviderStatus = errors.New("payment provider tidak mengembalikan status transaksi")

// ReconcilePayments mengecek status order yang masih PENDING ke payment provider
// lalu menyelaraskannya.
//
// Tanpa ini, order yang webhook-nya hilang akan menggantung di PENDING sampai cron
// kedaluwarsa membatalkannya, padahal uang pembeli sudah masuk.
// GetTransactionStatus sudah tersedia di provider tapi sebelumnya tidak pernah dipanggil.
func (s *service) ReconcilePayments() (int, error) {
	orders, err := s.orderRepo.FindPendingForReconciliation(200)
	if err != nil {
		return 0, err
	}

	updated := 0
	for _, o := range orders {
		resp, err := s.midtransProvider.GetTransactionStatus(o.ID)
		if err != nil || resp == nil || resp.TransactionStatus == "" {
			// Kegagalan satu order tidak boleh menghentikan sisa rekonsiliasi.
			log.Printf("reconcile order %s: gagal mengambil status dari provider: %v", o.ID, err)
			continue
		}

		// Payload dibentuk lalu dilewatkan HandleNotification supaya ada satu jalur
		// pemetaan status dan satu tempat aturan gross amount. Signature tidak
		// diverifikasi karena sumbernya adalah panggilan API terotorisasi, bukan webhook.
		payload := midtransProvider.NotificationPayload{
			OrderID:           o.ID,
			TransactionStatus: resp.TransactionStatus,
			FraudStatus:       "accept",
			GrossAmount:       strconv.FormatInt(int64(o.TotalAmount), 10),
		}

		if _, _, err := s.HandleNotification(payload); err != nil {
			log.Printf("reconcile order %s: gagal menerapkan status: %v", o.ID, err)
			continue
		}
		updated++
	}

	return updated, nil
}
