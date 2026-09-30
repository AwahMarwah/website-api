package task

import "github.com/hibiken/asynq"

// TypeCancelExpiredOrders membatalkan order PENDING yang sudah melewati masa berlaku
// pembayaran, sekaligus mengembalikan stoknya.
const TypeCancelExpiredOrders = "order:cancel_expired"

func NewCancelExpiredOrdersTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeCancelExpiredOrders, nil), nil
}

// NewReconcilePaymentsTask membuat task rekonsiliasi pembayaran. Dijadwalkan lebih
// jarang dari pembatalan kedaluwarsa supaya tidak membanjiri payment provider.
func NewReconcilePaymentsTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeReconcilePayments, nil), nil
}
