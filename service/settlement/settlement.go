package settlement

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	settlementModel "website-api/model/settlement"
	"website-api/repository/settlement"
)

type (
	// CommissionItem adalah satu baris order yang komisinya perlu dihitung.
	CommissionItem struct {
		OrderItemID string
		Subtotal    float64
	}

	IService interface {
		RecordCommission(tx *gorm.DB, merchantID, orderID string, rateBP int, items []CommissionItem) error
		RecordRefund(tx *gorm.DB, merchantID, orderID, note string, amount float64) error
		ListLedgers(merchantID string, reqQuery *settlementModel.ListLedgerReqQuery) ([]settlementModel.MerchantLedgerResponse, int64, int, error)
		Balance(merchantID string) (settlementModel.BalanceResponse, int, error)
	}

	service struct {
		settlementRepo settlement.IRepo
	}
)

func NewService(settlementRepo settlement.IRepo) IService {
	return &service{settlementRepo: settlementRepo}
}

// RecordCommission mencatat komisi marketplace untuk item-item milik satu merchant.
//
// Panggilannya idempoten terhadap OrderItemID karena setiap baris ledger item punya
// UNIQUE constraint. Kalau handler order dijalankan ulang, item yang sudah pernah
// dicatat dilewati dan saldo tidak berubah.
//
// amount pada ledger bernilai negatif: komisi adalah potongan dari pendapatan seller,
// bukan pemasukan. Dengan begitu saldo = SUM(amount) selalu langsung benar tanpa
// perlu membedakan arah per jenis baris.
func (s *service) RecordCommission(tx *gorm.DB, merchantID, orderID string, rateBP int, items []CommissionItem) error {
	if len(items) == 0 {
		return nil
	}

	repo := s.settlementRepo.WithTx(tx)
	ledgerID := uuid.NewString()

	balance, err := repo.GetBalance(merchantID)
	if err != nil {
		return err
	}

	totalCommission := 0.0
	claimed := make([]settlementModel.MerchantLedgerItem, 0, len(items))

	for _, item := range items {
		commission := CalculateCommission(item.Subtotal, rateBP)
		ok, err := repo.ClaimOrderItem(settlementModel.MerchantLedgerItem{
			ID:               uuid.NewString(),
			LedgerID:         ledgerID,
			OrderItemID:      item.OrderItemID,
			Subtotal:         item.Subtotal,
			CommissionRateBP: rateBP,
			Commission:       commission,
			CreatedAt:        time.Now(),
		})
		if err != nil {
			return err
		}
		// Item yang sudah pernah dicatat dilewati supaya tidak menggandakan saldo.
		if !ok {
			continue
		}
		claimed = append(claimed, settlementModel.MerchantLedgerItem{
			OrderItemID: item.OrderItemID,
			Subtotal:    item.Subtotal,
			Commission:  commission,
		})
		totalCommission += commission
	}

	if len(claimed) == 0 {
		return nil
	}

	note := "komisi marketplace"
	return repo.CreateLedger(settlementModel.MerchantLedger{
		ID:           ledgerID,
		MerchantID:   merchantID,
		OrderID:      &orderID,
		Type:         settlementModel.LedgerTypeCommission,
		Amount:       -totalCommission,
		BalanceAfter: balance - totalCommission,
		Note:         &note,
		CreatedAt:    time.Now(),
	})
}

// CalculateCommission menghitung komisi dari subtotal dalam basis points.
//
// Sengaja terkumpul di satu fungsi supaya seluruh tempat yang menghitung komisi
// memakai pembulatan yang sama. Saat tipe uang migrate ke integer, cukup ubah
// fungsi ini.
func CalculateCommission(subtotal float64, rateBP int) float64 {
	if rateBP <= 0 || subtotal <= 0 {
		return 0
	}
	commission := subtotal * float64(rateBP) / 10000
	return roundCents(commission)
}

func roundCents(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
