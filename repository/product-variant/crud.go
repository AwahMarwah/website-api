package product_variant

import (
	"gorm.io/gorm"

	product_variant "website-api/model/product-variant"
)

func (r *repo) Create(productVariant product_variant.ProductVariant) error {
	return r.db.Create(&productVariant).Error
}

func (r *repo) UpdateFields(id string, values map[string]any) error {
	return r.db.Model(&product_variant.ProductVariant{}).Where("id = ?", id).Updates(values).Error
}

// Restock mengembalikan stok tanpa mengecek batas atas. Dipakai saat pembatalan
// atau refund, di mana stok sebelumnya benar-benar sudah dikurangi.
func (r *repo) Restock(variantID string, qty int) error {
	return r.db.Model(&product_variant.ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock", gorm.Expr("stock + ?", qty)).Error
}

func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{db: tx}
}
