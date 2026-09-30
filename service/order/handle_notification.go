package order

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"website-api/common"
	"website-api/model/order"
	userModel "website-api/model/user"
	"website-api/task"
	midtransProvider "website-api/third-party/provider/midtrans"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// HandleNotification memproses notifikasi webhook dari Midtrans.
// Payload diverifikasi signature, lalu status order diperbarui.
//
// Notifikasi Midtrans bisa dikirim berulang dan beruntun (capture lalu settlement).
// Karena itu perpindahan status memakai UpdateStatusFrom yang bersyarat: kalau status
// sudah lebih maju, notifikasi lama diabaikan dan tidak memicu efek samping apa pun.
func (s *service) HandleNotification(payload midtransProvider.NotificationPayload) (order.NotificationResponse, int, error) {
	resData := order.NotificationResponse{
		OrderID:   payload.OrderID,
		Processed: false,
	}

	existingOrder, err := s.orderRepo.FindByID(payload.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resData.Status = "NOT_FOUND"
			return resData, http.StatusNotFound, fmt.Errorf("order not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil order: %w", err)
	}

	newStatus := payload.ResolveStatus()

	// jangan proses ulang order yang sudah final
	if isTerminalStatus(existingOrder.Status) {
		resData.Status = existingOrder.Status
		resData.Processed = true
		return resData, http.StatusOK, nil
	}

	// Verifikasi jumlah yang dibayar sama dengan total order. Tanpa ini, webhook
	// tervalidasi signature tapi dengan gross_amount buatan sendiri bisa mengubah status.
	gross, parseErr := parseGrossAmount(payload.GrossAmount)
	if parseErr != nil {
		return resData, http.StatusBadRequest, fmt.Errorf("invalid gross amount: %w", parseErr)
	}
	if gross != int64(existingOrder.TotalAmount) {
		resData.Status = "MISMATCH"
		return resData, http.StatusBadRequest, fmt.Errorf("gross amount mismatch with order total")
	}

	// Notifikasi refund tidak boleh menaikkan status order menjadi CANCELLED diam-diam;
	// itu harus lewat flow refund yang diaudit.
	if payload.IsRefund() {
		resData.Status = existingOrder.Status
		resData.Processed = true
		log.Printf("order %s menerima notifikasi refund, menunggu proses refund terotorisasi", existingOrder.ID)
		return resData, http.StatusOK, nil
	}

	if err := s.applyPaymentStatusChange(existingOrder, newStatus); err != nil {
		return resData, http.StatusInternalServerError, err
	}

	resData.Status = newStatus
	resData.Processed = true

	// Kirim email hanya saat transisi benar-benar terjadi, bukan saat notifikasi duplikat.
	if newStatus == common.OrderStatusPaid {
		s.enqueuePaymentSuccessEmail(existingOrder)
	}

	return resData, http.StatusOK, nil
}

func isTerminalStatus(status string) bool {
	switch status {
	case common.OrderStatusPaid, common.OrderStatusCompleted,
		common.OrderStatusCancelled, common.OrderStatusExpired, common.OrderStatusRefunded:
		return true
	}
	return false
}

// applyPaymentStatusChange menyimpan status pembayaran. Berbeda dari pembatalan,
// pembayaran tidak mengembalikan stok karena stok memang dialokasikan untuk pesanan ini.
func (s *service) applyPaymentStatusChange(existing order.Order, newStatus string) error {
	return s.txManager.Execute(func(tx *gorm.DB) error {
		txOrderRepo := s.orderRepo.WithTx(tx)

		affected, err := txOrderRepo.UpdateStatusFrom(existing.ID, existing.Status, newStatus)
		if err != nil {
			return fmt.Errorf("gagal memperbarui status order: %w", err)
		}
		if affected == 0 {
			// notifikasi duplikat, status sudah diubah request sebelumnya
			return nil
		}

		if err := txOrderRepo.UpdateStatusTimestamps(existing.ID, newStatus); err != nil {
			return fmt.Errorf("gagal memperbarui timestamp status: %w", err)
		}

		from := existing.Status
		note := "notifikasi pembayaran " + newStatus
		return txOrderRepo.CreateStatusHistory(order.OrderStatusHistory{
			OrderID:    existing.ID,
			FromStatus: &from,
			ToStatus:   newStatus,
			ActorRole:  optionalString("payment_gateway"),
			Note:       &note,
			CreatedAt:  time.Now(),
		})
	})
}

// enqueuePaymentSuccessEmail mengirim email invoice secara async melalui Asynq
func (s *service) enqueuePaymentSuccessEmail(order order.Order) {
	user, err := s.userRepo.Take([]string{"id", "name", "email"}, &userModel.User{Id: order.UserID})
	if err != nil {
		log.Printf("failed get user %s for payment email: %v", order.UserID, err)
		return
	}

	emailTask, err := task.NewPaymentSuccessTask(
		user.Name,
		user.Email,
		order.ID,
		order.TotalAmount,
	)
	if err != nil {
		log.Printf("failed create payment email task: %v", err)
		return
	}

	if _, err := s.queueClient.Enqueue(
		emailTask,
		asynq.MaxRetry(3),
		asynq.Queue("critical"),
	); err != nil {
		log.Printf("failed enqueue payment email: %v", err)
	}
}
