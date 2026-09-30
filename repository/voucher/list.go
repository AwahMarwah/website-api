package voucher

import voucherModel "website-api/model/voucher"

func (r *repo) List(reqQuery *voucherModel.ListVoucherReqQuery) (vouchers []voucherModel.Voucher, count int64, err error) {
	query := r.db.Model(&voucherModel.Voucher{})

	if reqQuery.Search != "" {
		query = query.Where("code ILIKE ?", "%"+reqQuery.Search+"%")
	}
	if reqQuery.IsActive != nil {
		query = query.Where("is_active = ?", *reqQuery.IsActive)
	}
	if reqQuery.MerchantID != "" {
		query = query.Where("merchant_id = ?", reqQuery.MerchantID)
	}

	if err = query.Count(&count).Error; err != nil {
		return nil, count, err
	}
	err = query.Limit(reqQuery.Limit).Offset(reqQuery.Offset).
		Order("created_at DESC").Find(&vouchers).Error
	return vouchers, count, err
}
