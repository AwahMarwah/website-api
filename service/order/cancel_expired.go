package order

import (
	"log"
	"net/http"
	"time"

	"website-api/common"
)

// CancelExpiredOrders membatalkan order yang masih PENDING namun sudah melewati expired_at,
// sekaligus mengembalikan stok. Dijalankan oleh scheduler Asynq setiap 5 menit.
func (s *service) CancelExpiredOrders() (int, error) {
	expiredOrders, err := s.orderRepo.FindExpiredPending(time.Now())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, o := range expiredOrders {
		// Transisi dari PENDING dijaga UpdateStatusFrom, jadi order yang keburu dibayar
		// lewat webhook di detik yang sama tidak akan ikut mengembalikan stok.
		if _, err := s.cancelAndRestoreStock(o.ID, common.OrderStatusExpired, systemActor(), "kedaluwarsa tanpa pembayaran"); err != nil {
			log.Printf("failed to expire order %s: %v", o.ID, err)
			continue
		}
	}

	return http.StatusOK, nil
}
