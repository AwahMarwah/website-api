package product_variant

import (
	"gorm.io/gorm/clause"

	product_variant "website-api/model/product-variant"
)

// FindByIDForUpdate mengambil variant dengan SELECT ... FOR UPDATE.
//
// Ini wajib dipakai di dalam transaksi checkout. Tanpa baris terkunci, dua request
// bersamaan untuk varian yang sama bisa sama-sama membaca stok cukup dan sama-sama
// mengurangi stok, sehingga terjual lebih banyak dari stok yang ada.
func (r *repo) FindByIDForUpdate(id string) (productVariant product_variant.ProductVariant, err error) {
	return productVariant, r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&productVariant).Error
}
