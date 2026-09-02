package order

import (
	"log"
	"net/http"
	"time"
)

// CancelExpiredOrders membatalkan order yang masih PENDING namun sudah melewati expired_at.
func (s *service) CancelExpiredOrders() (int, error) {
	expiredOrders, err := s.orderRepo.FindExpiredPending(time.Now())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, o := range expiredOrders {
		if err := s.orderRepo.UpdateStatus(o.ID, "EXPIRED"); err != nil {
			log.Printf("failed to expire order %s: %v", o.ID, err)
			continue
		}
	}

	return http.StatusOK, nil
}
