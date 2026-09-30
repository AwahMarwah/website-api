package settlement

import settlementModel "website-api/model/settlement"

// getBalanceSQL menjumlahkan kolom amount yang sudah bertanda.
const getBalanceSQL = `SELECT COALESCE(SUM(amount), 0) FROM merchant_ledgers WHERE merchant_id = ?`

func (r *repo) GetBalance(merchantID string) (float64, error) {
	var balance float64
	err := r.db.Raw(getBalanceSQL, merchantID).Scan(&balance).Error
	return balance, err
}

const sumByTypeSQL = `SELECT COALESCE(SUM(amount), 0) FROM merchant_ledgers WHERE merchant_id = ? AND type = ?`

func (r *repo) SumCommission(merchantID string) (float64, error) {
	var total float64
	err := r.db.Raw(sumByTypeSQL, merchantID, settlementModel.LedgerTypeCommission).Scan(&total).Error
	// Dikembalikan positif supaya presentation layer tidak perlu membalik tanda lagi.
	return -total, err
}

func (r *repo) SumPayout(merchantID string) (float64, error) {
	var total float64
	err := r.db.Raw(sumByTypeSQL, merchantID, settlementModel.LedgerTypePayout).Scan(&total).Error
	return -total, err
}
