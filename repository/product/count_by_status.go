package product

import (
	productModel "website-api/model/product"

	"gorm.io/gorm"
)

// CountByStatus mengembalikan jumlah produk per status untuk satu merchant,
// dipakai ringkasan dashboard seller.
func (r *repo) CountByStatus(merchantID string) *gorm.DB {
	return r.db.Model(&productModel.Product{}).
		Select("status, COUNT(*) AS total").
		Where("merchant_id = ?", merchantID).
		Group("status")
}
