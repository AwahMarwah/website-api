package product

import (
	"database/sql"

	productModel "website-api/model/product"

	"gorm.io/gorm"
)

// findCheckoutSnapshotSQL mengambil nama produk, merchant pemilik, dan gambar utama
// dalam satu round trip. Subquery memilih gambar dengan sort_order terkecil supaya
// hasilnya deterministik saat belum ada gambar yang ditandai primary.
const findCheckoutSnapshotSQL = `
SELECT p.id            AS product_id,
       p.name          AS product_name,
       COALESCE(p.merchant_id, '') AS merchant_id,
       COALESCE((
           SELECT pi.image_url
           FROM product_images pi
           WHERE pi.product_id = p.id
           ORDER BY pi.is_primary DESC, pi.sort_order ASC, pi.created_at ASC
           LIMIT 1
       ), '') AS image_url
FROM products p
WHERE p.id = ?
`

func (r *repo) FindCheckoutSnapshot(productID string) (productModel.CheckoutSnapshot, error) {
	var snapshot productModel.CheckoutSnapshot
	err := r.db.Raw(findCheckoutSnapshotSQL, productID).Scan(&snapshot).Error
	if err == sql.ErrNoRows {
		return snapshot, gorm.ErrRecordNotFound
	}
	return snapshot, err
}
