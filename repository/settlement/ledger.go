package settlement

import (
	"gorm.io/gorm/clause"

	settlementModel "website-api/model/settlement"
)

// clauseOnConflictDoNothing membuat INSERT diabaikan diam-diam kalau OrderItemID
// sudah ada, sehingga komisi tidak pernah terhitung dua kali.
var clauseOnConflictDoNothing = clause.OnConflict{DoNothing: true}

func (r *repo) ListLedgers(merchantID, ledgerType string, limit, offset int) ([]settlementModel.MerchantLedger, int64, error) {
	query := r.db.Model(&settlementModel.MerchantLedger{}).Where("merchant_id = ?", merchantID)
	if ledgerType != "" {
		query = query.Where("type = ?", ledgerType)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, count, err
	}

	var ledgers []settlementModel.MerchantLedger
	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&ledgers).Error
	return ledgers, count, err
}

func (r *repo) CreateLedger(ledger settlementModel.MerchantLedger) error {
	return r.db.Create(&ledger).Error
}

func (r *repo) CreateLedgerItem(item settlementModel.MerchantLedgerItem) error {
	return r.db.Create(&item).Error
}

// ClaimOrderItem mencoba mencadangkan satu item order untuk komisi.
// Mengembalikan false kalau OrderItemID sudah pernah dicatat.
func (r *repo) ClaimOrderItem(item settlementModel.MerchantLedgerItem) (bool, error) {
	result := r.db.Clauses(clauseOnConflictDoNothing).Create(&item)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}
