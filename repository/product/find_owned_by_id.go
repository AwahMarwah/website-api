package product

import (
	"database/sql"

	productModel "website-api/model/product"

	"gorm.io/gorm"
)

// FindOwnedByID mengambil produk hanya jika produk itu milik merchant yang diberikan.
// Dipakai sebagai guard ownership seller: query menolak di level database sehingga
// tidak ada jalur kode yang bisa hanya karena lupa mengecek.
func (r *repo) FindOwnedByID(id, merchantID string) (product productModel.Product, err error) {
	err = r.db.Where("id = ? AND merchant_id = ?", id, merchantID).First(&product).Error
	if err == sql.ErrNoRows {
		return product, gorm.ErrRecordNotFound
	}
	return product, err
}

func (r *repo) CountOwnedByStatus(merchantID, status string) (int64, error) {
	var count int64
	query := r.db.Model(&productModel.Product{}).Where("merchant_id = ?", merchantID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return count, query.Count(&count).Error
}
