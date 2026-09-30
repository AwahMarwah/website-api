package settlement

import (
	"database/sql"

	"gorm.io/gorm"

	settlementModel "website-api/model/settlement"
)

func (r *repo) CreatePayout(payout settlementModel.Payout) error {
	return r.db.Create(&payout).Error
}

func (r *repo) FindPayoutByID(id string) (settlementModel.Payout, error) {
	var payout settlementModel.Payout
	err := r.db.Where("id = ?", id).First(&payout).Error
	if err == sql.ErrNoRows {
		return payout, gorm.ErrRecordNotFound
	}
	return payout, err
}

func (r *repo) ListPayouts(reqQuery *settlementModel.ListPayoutReqQuery) ([]settlementModel.Payout, int64, error) {
	query := r.db.Model(&settlementModel.Payout{})

	if reqQuery.Status != "" {
		query = query.Where("status = ?", reqQuery.Status)
	}
	if reqQuery.MerchantID != "" {
		query = query.Where("merchant_id = ?", reqQuery.MerchantID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, count, err
	}

	var payouts []settlementModel.Payout
	err := query.Limit(reqQuery.Limit).Offset(reqQuery.Offset).
		Order("created_at DESC").Find(&payouts).Error
	return payouts, count, err
}

func (r *repo) UpdatePayout(id string, values map[string]any) error {
	return r.db.Model(&settlementModel.Payout{}).Where("id = ?", id).Updates(values).Error
}
