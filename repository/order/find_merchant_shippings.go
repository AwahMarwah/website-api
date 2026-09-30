package order

import "website-api/model/order"

func (r *repo) FindMerchantShippingsByOrder(orderID string) ([]order.OrderMerchantShipping, error) {
	var shippings []order.OrderMerchantShipping
	err := r.db.
		Preload("Merchant").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&shippings).Error
	return shippings, err
}
