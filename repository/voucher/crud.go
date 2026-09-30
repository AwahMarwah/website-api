package voucher

import (
	"gorm.io/gorm"

	voucherModel "website-api/model/voucher"
)

func (r *repo) Create(v voucherModel.Voucher) error {
	return r.db.Create(&v).Error
}

func (r *repo) Update(id string, values map[string]any) error {
	return r.db.Model(&voucherModel.Voucher{}).Where("id = ?", id).Updates(values).Error
}

func (r *repo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&voucherModel.Voucher{}).Error
}

// ConsumeQuota menambah used_count hanya selama kuota belum habis.
// Syaratnya ada di dalam statement UPDATE, bukan dibaca dulu lalu ditulis,
// sehingga dua checkout bersamaan tidak bisa sama-sama melewati batas kuota.
func (r *repo) ConsumeQuota(id string) (int64, error) {
	result := r.db.Model(&voucherModel.Voucher{}).
		Where("id = ? AND (quota IS NULL OR used_count < quota)", id).
		Update("used_count", gorm.Expr("used_count + 1"))
	return result.RowsAffected, result.Error
}

// ReleaseQuota mengembalikan satu kuota, dipakai saat order gagal setelah voucher sempat diklaim.
func (r *repo) ReleaseQuota(id string) error {
	return r.db.Model(&voucherModel.Voucher{}).
		Where("id = ? AND used_count > 0", id).
		Update("used_count", gorm.Expr("used_count - 1")).Error
}

func (r *repo) CountUserRedemptions(voucherID, userID string) (int64, error) {
	var count int64
	err := r.db.Model(&voucherModel.VoucherRedemption{}).
		Where("voucher_id = ? AND user_id = ?", voucherID, userID).
		Count(&count).Error
	return count, err
}

func (r *repo) CreateRedemption(redemption voucherModel.VoucherRedemption) error {
	return r.db.Create(&redemption).Error
}

func (r *repo) RebindProducts(voucherID string, productIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("voucher_id = ?", voucherID).Delete(&voucherModel.ProductVoucher{}).Error; err != nil {
			return err
		}
		if len(productIDs) == 0 {
			return nil
		}
		rows := make([]voucherModel.ProductVoucher, 0, len(productIDs))
		for _, productID := range productIDs {
			rows = append(rows, voucherModel.ProductVoucher{VoucherID: voucherID, ProductID: productID})
		}
		return tx.Create(&rows).Error
	})
}

func (r *repo) FindProductIDs(voucherID string) ([]string, error) {
	var rows []voucherModel.ProductVoucher
	if err := r.db.Where("voucher_id = ?", voucherID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ProductID)
	}
	return ids, nil
}
