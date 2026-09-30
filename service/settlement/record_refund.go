package settlement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	settlementModel "website-api/model/settlement"
)

// RecordRefund mengembalikan komisi yang sebelumnya sudah dipotong ke merchant ketika
// order dibatalkan atau di-refund.
//
// Tanda amount dibalik: komisi dicatat negatif, refund dicatat positif, sehingga
// satu order yang dibayar lalu dibatalkan tidak meninggalkan saldo positif palsu.
func (s *service) RecordRefund(tx *gorm.DB, merchantID, orderID, note string, amount float64) error {
	if amount <= 0 {
		return nil
	}

	repo := s.settlementRepo.WithTx(tx)

	balance, err := repo.GetBalance(merchantID)
	if err != nil {
		return err
	}

	return repo.CreateLedger(settlementModel.MerchantLedger{
		ID:           uuid.NewString(),
		MerchantID:   merchantID,
		OrderID:      &orderID,
		Type:         settlementModel.LedgerTypeRefund,
		Amount:       amount,
		BalanceAfter: balance + amount,
		Note:         optionalString(note),
		CreatedAt:    time.Now(),
	})
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
