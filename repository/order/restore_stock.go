package order

// restoreStockForOrderSQL mengembalikan stok seluruh item dalam satu order.
// Ditulis sebagai satu statement agar tidak ada celah antara read dan write,
// dan tidak ada window di mana order sudah CANCELLED tapi stok belum kembali.
const restoreStockForOrderSQL = `
UPDATE product_variants pv
SET stock = pv.stock + agg.total_qty
FROM (SELECT product_variant_id, SUM(qty) AS total_qty
      FROM order_items
      WHERE order_id = ?
      GROUP BY product_variant_id) AS agg
WHERE pv.id = agg.product_variant_id
`

// RestoreStockForOrder mengembalikan stok untuk semua item pada order terkait.
// Mengembalikan jumlah variant yang berhasil di-restock.
func (r *repo) RestoreStockForOrder(orderID string) (int64, error) {
	result := r.db.Exec(restoreStockForOrderSQL, orderID)
	return result.RowsAffected, result.Error
}
