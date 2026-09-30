package cart

import (
	"gorm.io/gorm"

	"website-api/model/cart"
)

// upsertCartItemSQL menambahkan qty ke baris yang sudah ada untuk pasangan
// (user_id, product_variant_id). Tabel punya UNIQUE constraint pada kedua kolom
// itu, jadi INSERT biasa akan gagal dengan 500 kalau varian sama ditambahkan lagi.
const upsertCartItemSQL = `
INSERT INTO cart_items (id, user_id, product_variant_id, qty, created_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT (user_id, product_variant_id)
DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty
`

// Upsert menambah qty bila varian sudah ada di keranjang, atau membuat baris baru.
func (r *repo) Upsert(item *cart.CartItem) error {
	return r.db.Exec(upsertCartItemSQL, item.ID, item.UserID, item.ProductVariantID, item.Qty).Error
}

// SetQty menetapkan qty absolut (dipakai saat pembeli mengubah jumlah di keranjang).
func (r *repo) SetQty(userID, variantID string, qty int) error {
	return r.db.Model(&cart.CartItem{}).
		Where("user_id = ? AND product_variant_id = ?", userID, variantID).
		Update("qty", qty).Error
}

// WithTx mengikat repo ke transaksi yang sedang berjalan.
func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{db: tx}
}
