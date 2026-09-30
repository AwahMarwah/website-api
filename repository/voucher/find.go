package voucher

import (
	"database/sql"

	"gorm.io/gorm"

	voucherModel "website-api/model/voucher"
)

func (r *repo) FindByID(id string) (voucherModel.Voucher, error) {
	var v voucherModel.Voucher
	err := r.db.Where("id = ?", id).First(&v).Error
	if err == sql.ErrNoRows {
		return v, gorm.ErrRecordNotFound
	}
	return v, err
}

// FindByCode mencari voucher dengan case-insensitive supaya pembeli tidak gagal
// hanya karena kode diketik dengan huruf kecil.
func (r *repo) FindByCode(code string) (voucherModel.Voucher, error) {
	var v voucherModel.Voucher
	err := r.db.Where("UPPER(code) = ?", code).First(&v).Error
	if err == sql.ErrNoRows {
		return v, gorm.ErrRecordNotFound
	}
	return v, err
}
