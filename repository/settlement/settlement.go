package settlement

import (
	"gorm.io/gorm"

	settlementModel "website-api/model/settlement"
)

type (
	IRepo interface {
		// GetBalance menjumlahkan amount dari ledger. Nilai bertanda: komisi dan payout
		// negatif, refund positif.
		GetBalance(merchantID string) (float64, error)
		SumCommission(merchantID string) (float64, error)
		SumPayout(merchantID string) (float64, error)
		ListLedgers(merchantID, ledgerType string, limit, offset int) ([]settlementModel.MerchantLedger, int64, error)
		CreateLedger(ledger settlementModel.MerchantLedger) error
		CreateLedgerItem(item settlementModel.MerchantLedgerItem) error
		// ClaimOrderItem mengembalikan false kalau item order sudah pernah dicatat
		// komisinya. Berbasis UNIQUE constraint pada order_item_id.
		ClaimOrderItem(item settlementModel.MerchantLedgerItem) (bool, error)

		CreatePayout(payout settlementModel.Payout) error
		FindPayoutByID(id string) (settlementModel.Payout, error)
		FindPayoutByIDForUpdate(id string) (settlementModel.Payout, error)
		ListPayouts(reqQuery *settlementModel.ListPayoutReqQuery) ([]settlementModel.Payout, int64, error)
		UpdatePayout(id string, values map[string]any) error
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{db: tx}
}
