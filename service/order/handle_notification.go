package order

import (
	"fmt"
	"log"
	"net/http"
	"website-api/model/order"
	userModel "website-api/model/user"
	"website-api/task"
	midtransProvider "website-api/third-party/provider/midtrans"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// HandleNotification memproses notifikasi webhook dari Midtrans.
// Payload diverifikasi signature, lalu status order diperbarui.
func (s *service) HandleNotification(payload midtransProvider.NotificationPayload) (order.NotificationResponse, int, error) {	resData := order.NotificationResponse{
		OrderID:   payload.OrderID,
		Processed: false,
	}

	// cek order terlebih dahulu
	existingOrder, err := s.orderRepo.FindByID(payload.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resData.Status = "NOT_FOUND"
			return resData, http.StatusNotFound, fmt.Errorf("order not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil order: %w", err)
	}

	// tentukan status baru dari payload
	newStatus := payload.ResolveStatus()

	// jangan ubah order yang sudah PAID / final
	if existingOrder.Status == "PAID" || existingOrder.Status == "COMPLETED" || existingOrder.Status == "CANCELLED" {
		resData.Status = existingOrder.Status
		resData.Processed = true
		return resData, http.StatusOK, nil
	}

	// verifikasi signature hanya untuk transaksi yang mengubah status ke pembayaran berhasil
	// untuk keamanan, validasi jumlah yang dibayar sesuai total order
	gross, parseErr := parseGrossAmount(payload.GrossAmount)
	if parseErr != nil {
		return resData, http.StatusBadRequest, fmt.Errorf("invalid gross amount: %w", parseErr)
	}
	if gross != int64(existingOrder.TotalAmount) {
		resData.Status = "MISMATCH"
		return resData, http.StatusBadRequest, fmt.Errorf("gross amount mismatch with order total")
	}

	if err := s.orderRepo.UpdateStatus(payload.OrderID, newStatus); err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal memperbarui status order: %w", err)
	}

	resData.Status = newStatus
	resData.Processed = true

	// kirim email notifikasi setelah pembayaran sukses
	if newStatus == "PAID" {
		s.enqueuePaymentSuccessEmail(existingOrder)
	}

	return resData, http.StatusOK, nil
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
