package product

import (
	productModel "website-api/model/product"

	"gorm.io/gorm"
)

// FindDetailCollections mengambil product beserta koleksi relasinya (gambar)
// lewat mixed strategy: agregat tetap di FindDetailByID, koleksi row di-preload.
func (r *repo) FindDetailCollections(id string) (resData productModel.Product, err error) {
	return resData, r.db.
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Where("products.id = ?", id).
		First(&resData).Error
}