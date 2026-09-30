package settlement

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"website-api/database/transaction"
	"website-api/library/pagination"
	settlementModel "website-api/model/settlement"
	"website-api/repository/settlement"
)

type IPayoutService interface {
	// RequestPayout membuat pengajuan pencairan dana oleh seller.
	RequestPayout(merchantID, requestedBy string, req *settlementModel.CreatePayoutReq) (string, int, error)
	// ApprovePayout menyetujui pencairan dan mencatatnya ke ledger.
	ApprovePayout(payoutID, approverID string, req *settlementModel.ApprovePayoutReq) (int, error)
	// MarkPaid menandai pencairan sudah ditransfer ke rekening seller.
	MarkPaid(payoutID, approverID string) (int, error)
	RejectPayout(payoutID, approverID, reason string) (int, error)
	ListPayouts(reqQuery *settlementModel.ListPayoutReqQuery) ([]settlementModel.PayoutResponse, int64, int, error)
}

type payoutService struct {
	settlementRepo settlement.IRepo
	txManager      transaction.ITransactionManager
}

func NewPayoutService(settlementRepo settlement.IRepo, txManager transaction.ITransactionManager) IPayoutService {
	return &payoutService{settlementRepo: settlementRepo, txManager: txManager}
}

// RequestPayout membuat pengajuan pencairan.
//
// Nominal tidak langsung memotong saldo; saldo baru berkurang saat payout disetujui.
// Kalau dipotong sejak pengajuan, pengajuan yang akhirnya dibatalkan akan tetap
// mengurangi saldo seller.
func (s *payoutService) RequestPayout(merchantID, requestedBy string, req *settlementModel.CreatePayoutReq) (string, int, error) {
	balance, err := s.settlementRepo.GetBalance(merchantID)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mengambil saldo: %w", err)
	}
	if req.Amount > balance {
		return "", http.StatusBadRequest, fmt.Errorf("nominal pencairan melebihi saldo tersedia")
	}

	payout := settlementModel.Payout{
		ID:          uuid.NewString(),
		MerchantID:  merchantID,
		Amount:      req.Amount,
		Status:      settlementModel.PayoutStatusPending,
		BankAccount: &req.BankAccount,
		Note:        optionalString(req.Note),
		RequestedBy: &requestedBy,
		CreatedAt:   time.Now(),
	}
	if err := s.settlementRepo.CreatePayout(payout); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal membuat pengajuan payout: %w", err)
	}
	return payout.ID, http.StatusCreated, nil
}

// ApprovePayout mencatat pencairan ke ledger sehingga saldo seller berkurang.
//
// Baris payout dikunci dengan FOR UPDATE supaya dua approver yang menekan tombol
// bersamaan tidak menghasilkan dua kali pencairan.
func (s *payoutService) ApprovePayout(payoutID, approverID string, req *settlementModel.ApprovePayoutReq) (int, error) {
	err := s.txManager.Execute(func(tx *gorm.DB) error {
		repo := s.settlementRepo.WithTx(tx)

		current, err := repo.FindPayoutByIDForUpdate(payoutID)
		if err != nil {
			return err
		}
		if current.Status != settlementModel.PayoutStatusPending {
			return &payoutConflictError{msg: fmt.Sprintf("payout sudah diproses dengan status %s", current.Status)}
		}

		balance, err := repo.GetBalance(current.MerchantID)
		if err != nil {
			return err
		}
		if current.Amount > balance {
			return fmt.Errorf("nominal pencairan melebihi saldo tersedia")
		}

		ledgerID := uuid.NewString()
		note := "pencairan dana ke seller"
		if err := repo.CreateLedger(settlementModel.MerchantLedger{
			ID:           ledgerID,
			MerchantID:   current.MerchantID,
			Type:         settlementModel.LedgerTypePayout,
			Amount:       -current.Amount,
			BalanceAfter: balance - current.Amount,
			Note:         &note,
			CreatedAt:    time.Now(),
		}); err != nil {
			return err
		}

		return repo.UpdatePayout(payoutID, map[string]any{
			"status":      settlementModel.PayoutStatusApproved,
			"approved_by": &approverID,
			"ledger_id":   &ledgerID,
			"note":        optionalString(req.Note),
			"updated_at":  time.Now(),
		})
	})
	if err != nil {
		if conflict, ok := err.(*payoutConflictError); ok {
			return http.StatusConflict, conflict
		}
		if verr := asValidationError(err); verr != nil {
			return http.StatusBadRequest, verr
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal menyetujui payout: %w", err)
	}
	return http.StatusOK, nil
}

func (s *payoutService) MarkPaid(payoutID, approverID string) (int, error) {
	payout, err := s.settlementRepo.FindPayoutByID(payoutID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("payout not found")
	}
	if payout.Status != settlementModel.PayoutStatusApproved {
		return http.StatusConflict, fmt.Errorf("payout harus disetujui sebelum ditandai lunas")
	}

	if err := s.settlementRepo.UpdatePayout(payoutID, map[string]any{
		"status":      settlementModel.PayoutStatusPaid,
		"approved_by": &approverID,
		"updated_at":  time.Now(),
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menandai payout lunas: %w", err)
	}
	return http.StatusOK, nil
}

// RejectPayout menolak pengajuan. Saldo tidak terpengaruh karena pencairan baru
// dicatat ke ledger saat disetujui.
func (s *payoutService) RejectPayout(payoutID, approverID, reason string) (int, error) {
	payout, err := s.settlementRepo.FindPayoutByID(payoutID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("payout not found")
	}
	if payout.Status == settlementModel.PayoutStatusPaid {
		return http.StatusConflict, fmt.Errorf("payout yang sudah lunas tidak bisa ditolak")
	}

	if err := s.settlementRepo.UpdatePayout(payoutID, map[string]any{
		"status":      settlementModel.PayoutStatusRejected,
		"approved_by": &approverID,
		"note":        optionalString(reason),
		"updated_at":  time.Now(),
	}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menolak payout: %w", err)
	}
	return http.StatusOK, nil
}

func (s *payoutService) ListPayouts(reqQuery *settlementModel.ListPayoutReqQuery) ([]settlementModel.PayoutResponse, int64, int, error) {
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	payouts, count, err := s.settlementRepo.ListPayouts(reqQuery)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar payout: %w", err)
	}

	res := make([]settlementModel.PayoutResponse, 0, len(payouts))
	for _, p := range payouts {
		res = append(res, settlementModel.PayoutResponse{
			ID:          p.ID,
			MerchantID:  p.MerchantID,
			Amount:      p.Amount,
			Status:      p.Status,
			BankAccount: p.BankAccount,
			Note:        p.Note,
			RequestedBy: p.RequestedBy,
			ApprovedBy:  p.ApprovedBy,
			CreatedAt:   p.CreatedAt,
		})
	}
	return res, count, http.StatusOK, nil
}

// payoutConflictError menandai konflik state (payout sudah diproses) supaya service
// bisa mengembalikan 409, bukan 500.
type payoutConflictError struct{ msg string }

func (e *payoutConflictError) Error() string { return e.msg }

func asValidationError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errNotFound
	}
	return nil
}

var errNotFound = &payoutConflictError{msg: "payout not found"}
