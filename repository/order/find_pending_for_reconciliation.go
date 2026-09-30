package order

import (
	"time"

	"website-api/common"
	"website-api/model/order"
)

// FindPendingForReconciliation mengembalikan order yang masih menunggu pembayaran.
//
// Hanya order yang sudah lewat masa tenggang windows-nya yang diambil: order yang baru
// dibuat masih punya waktu untuk webhook-nya sampai, dan memanggil provider untuk
// setiap order baru akan membuang kuota API.
func (r *repo) FindPendingForReconciliation(limit int) ([]order.Order, error) {
	var orders []order.Order
	cutoff := time.Now().Add(-1 * time.Minute)

	err := r.db.
		Where("status = ? AND created_at <= ?", common.OrderStatusPending, cutoff).
		Limit(limit).
		Order("created_at ASC").
		Find(&orders).Error
	return orders, err
}
