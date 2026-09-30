package order

import (
	"gorm.io/gorm"

	"website-api/model/order"
)

// FindByMerchantID mengembalikan order yang minimal memuat satu item milik merchant.
// Query lewat order_items.merchant_id (yang sudah dibekukan saat checkout), bukan lewat
// join ke produk, supaya hasilnya tidak berubah ketika seller memindahkan produk.
func (r *repo) FindByMerchantID(merchantID, status string, limit, offset int) (orders []order.Order, count int64, err error) {
	query := r.db.Model(&order.Order{}).
		Where("id IN (?)", r.db.Model(&order.OrderItem{}).
			Select("order_id").
			Where("merchant_id = ?", merchantID))

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err = query.Count(&count).Error; err != nil {
		return nil, count, err
	}

	err = query.
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Where("merchant_id = ?", merchantID).
				Preload("ProductVariant.Product").
				Preload("Merchant")
		}).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, count, err
}

// FindItemsByOrderAndMerchant mengembalikan hanya item milik merchant tersebut.
// Seller tidak boleh melihat item milik seller lain yang berada di order yang sama.
func (r *repo) FindItemsByOrderAndMerchant(orderID, merchantID string) ([]order.OrderItem, error) {
	var items []order.OrderItem
	err := r.db.
		Preload("ProductVariant.Product").
		Where("order_id = ? AND merchant_id = ?", orderID, merchantID).
		Find(&items).Error
	return items, err
}

// MerchantOwnsOrder menjawab apakah order memuat barang dari merchant tersebut,
// dipakai sebagai guard ownership di layer service.
func (r *repo) MerchantOwnsOrder(orderID, merchantID string) (bool, error) {
	var count int64
	err := r.db.Model(&order.OrderItem{}).
		Where("order_id = ? AND merchant_id = ?", orderID, merchantID).
		Count(&count).Error
	return count > 0, err
}
