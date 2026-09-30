package order

import (
	"gorm.io/gorm"

	"website-api/common"
	"website-api/model/order"
)

// UpdateStatusFrom melakukan transisi status secara atomik dengan syarat status saat ini.
// Melindungi terhadap race: dua request pembatalan bersamaan tidak akan sama-sama
// mengembalikan stok. RowsAffected == 0 berarti order sudah diproses request lain.
func (r *repo) UpdateStatusFrom(id, fromStatus, toStatus string) (int64, error) {
	result := r.db.Model(&order.Order{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Update("status", toStatus)
	return result.RowsAffected, result.Error
}

// UpdateStatusTimestamps mengisi kolom waktu yang sesuai dengan status tujuan,
// sehingga histori level timeline tidak perlu diturunkan dari created_at.
func (r *repo) UpdateStatusTimestamps(id, status string) error {
	column := statusTimestampColumn(status)
	if column == "" {
		return nil
	}
	return r.db.Model(&order.Order{}).
		Where("id = ?", id).
		Update(column, gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func statusTimestampColumn(status string) string {
	switch status {
	case common.OrderStatusPaid:
		return "paid_at"
	case common.OrderStatusShipped:
		return "shipped_at"
	case common.OrderStatusCompleted:
		return "completed_at"
	case common.OrderStatusCancelled, common.OrderStatusExpired:
		return "cancelled_at"
	}
	return ""
}
