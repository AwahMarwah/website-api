package refund

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"website-api/common"
	"website-api/database/transaction"
	refundModel "website-api/model/refund"
	"website-api/repository/order"
	refundRepo "website-api/repository/refund"
	settlementRepo "website-api/repository/settlement"
	settlementService "website-api/service/settlement"
	midtransProvider "website-api/third-party/provider/midtrans"
)

type (
	// PaymentProvider adalah bagian dari provider pembayaran yang dipakai refund.
	// Sengaja interface sempit supaya service tidak terikat ke Midtrans.
	PaymentProvider interface {
		Refund(orderID, transactionStatus, reason string, amount int64) (*midtransProvider.RefundResponse, error)
	}

	IService interface {
		// Request mendaftarkan permintaan refund dari pembeli. Statusnya PENDING;
		// dana belum berpindah sampai disetujui finance.
		Request(orderID, userID string, req *refundModel.CreateRefundReq) (refundID string, statusCode int, err error)
		Approve(refundID, approverID string) (int, error)
		Reject(refundID, approverID, reason string) (int, error)
		Process(refundID string) (int, error)
		List(reqQuery *refundModel.ListRefundReqQuery) ([]refundModel.RefundResponse, int64, int, error)
		Detail(refundID string) (refundModel.RefundResponse, int, error)
	}

	service struct {
		refundRepo        refundRepo.IRepo
		orderRepo         order.IRepo
		settlementRepo    settlementRepo.IRepo
		settlementService settlementService.IService
		txManager         transaction.ITransactionManager
		provider          PaymentProvider
	}
)

func NewService(
	refundRepo refundRepo.IRepo,
	orderRepo order.IRepo,
	settlementRepo settlementRepo.IRepo,
	settlementService settlementService.IService,
	txManager transaction.ITransactionManager,
	provider PaymentProvider,
) IService {
	return &service{
		refundRepo:        refundRepo,
		orderRepo:         orderRepo,
		settlementRepo:    settlementRepo,
		settlementService: settlementService,
		txManager:         txManager,
		provider:          provider,
	}
}

// Request mendaftarkan permintaan refund.
//
// Nominal diverifikasi terhadap sisa yang belum direfund, bukan hanya terhadap total
// order: kalau sudah ada refund sebelumnya, permintaan kedua tidak boleh melebihi sisanya.
func (s *service) Request(orderID, userID string, req *refundModel.CreateRefundReq) (string, int, error) {
	o, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return "", http.StatusNotFound, fmt.Errorf("order not found")
	}
	if o.UserID != userID {
		return "", http.StatusForbidden, fmt.Errorf("forbidden")
	}
	if o.Status != common.OrderStatusPaid && o.Status != common.OrderStatusProcessing &&
		o.Status != common.OrderStatusShipped && o.Status != common.OrderStatusCompleted {
		return "", http.StatusBadRequest,
			fmt.Errorf("hanya order yang sudah dibayar yang bisa direfund")
	}

	alreadyRefunded, err := s.refundRepo.SumCompletedByOrder(orderID)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal menghitung total refund: %w", err)
	}
	remaining := o.TotalAmount - alreadyRefunded
	if req.Amount > remaining {
		return "", http.StatusBadRequest,
			fmt.Errorf("nominal refund melebihi sisa yang dapat direfund (%.0f)", remaining)
	}

	if pending, err := s.refundRepo.FindPendingByOrderID(orderID); err == nil && len(pending) > 0 {
		return "", http.StatusConflict, fmt.Errorf("sudah ada permintaan refund yang sedang diproses")
	}

	v := refundModel.Refund{
		ID:           uuid.NewString(),
		OrderID:      orderID,
		RequestedBy:  userID,
		Amount:       req.Amount,
		Reason:       &req.Reason,
		Status:       refundModel.StatusPending,
		Provider:     midtransProvider.ProviderName,
		Restocked:    false,
		CreatedAt:    time.Now(),
	}
	if err := s.refundRepo.Create(v); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal membuat permintaan refund: %w", err)
	}
	return v.ID, http.StatusCreated, nil
}
