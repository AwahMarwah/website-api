package order

import "website-api/model/order"

func (r *repo) CreateStatusHistory(history order.OrderStatusHistory) error {
	return r.db.Create(&history).Error
}

func (r *repo) FindStatusHistories(orderID string) ([]order.OrderStatusHistory, error) {
	var histories []order.OrderStatusHistory
	err := r.db.Where("order_id = ?", orderID).Order("created_at ASC").Find(&histories).Error
	return histories, err
}
