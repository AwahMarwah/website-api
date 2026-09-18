package order

import "website-api/model/order"

// HasCompletedOrderForProduct cek apakah user pernah punya order COMPLETED yang berisi produk tsb.
func (r *repo) HasCompletedOrderForProduct(userID, productID string) (bool, error) {
	var count int64
	err := r.db.Model(&order.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN product_variants pv ON pv.id = order_items.product_variant_id").
		Where("orders.user_id = ? AND orders.status = ? AND pv.product_id = ?",
			userID, "COMPLETED", productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}