package order

import "website-api/model/order"

func (r *repo) FindByIDWithItems(id string) (orderModel order.Order, items []order.OrderItem, err error) {
	err = r.db.
		Preload("Items.ProductVariant.Product").
		Where("id = ?", id).
		First(&orderModel).Error
	if err != nil {
		return orderModel, nil, err
	}
	return orderModel, orderModel.Items, nil
}

func (r *repo) FindByUserID(userID, status string, limit, offset int) (orders []order.Order, count int64, err error) {
	query := r.db.Model(&order.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err = query.Count(&count).Error; err != nil {
		return nil, count, err
	}
	err = query.
		Preload("Items.ProductVariant.Product").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, count, err
}

func (r *repo) FindAll(status string, limit, offset int) (orders []order.Order, count int64, err error) {
	query := r.db.Model(&order.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err = query.Count(&count).Error; err != nil {
		return nil, count, err
	}
	err = query.
		Preload("Items.ProductVariant.Product").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, count, err
}
