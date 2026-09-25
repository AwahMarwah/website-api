package order

import (
	"time"
	"website-api/model/order"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		FindByID(id string) (order.Order, error)
		FindByIDWithItems(id string) (order.Order, []order.OrderItem, error)
		FindByUserID(userID string, status string, limit, offset int) ([]order.Order, int64, error)
		FindAll(status string, limit, offset int) ([]order.Order, int64, error)
		Update(order.Order) error
		UpdateStatus(id, status string) error
		UpdatePaymentInfo(id string, values map[string]interface{}) error
		CreateMerchantShipping(shipping *order.OrderMerchantShipping) error
		HasCompletedOrderForProduct(userID, productID string) (bool, error)
		FindExpiredPending(now time.Time) ([]order.Order, error)
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{db: tx}
}

func (r *repo) FindByID(id string) (orderModel order.Order, err error) {
	return orderModel, r.db.Where("id = ?", id).First(&orderModel).Error
}

func (r *repo) Update(orderModel order.Order) error {
	return r.db.Save(&orderModel).Error
}

func (r *repo) UpdateStatus(id, status string) error {
	return r.db.Model(&order.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *repo) UpdatePaymentInfo(id string, values map[string]interface{}) error {
	return r.db.Model(&order.Order{}).Where("id = ?", id).Updates(values).Error
}

func (r *repo) CreateMerchantShipping(shipping *order.OrderMerchantShipping) error {
	return r.db.Create(shipping).Error
}

func (r *repo) FindExpiredPending(now time.Time) ([]order.Order, error) {
	var orders []order.Order
	err := r.db.
		Where("status = ? AND expired_at IS NOT NULL AND expired_at < ?", "PENDING", now).
		Find(&orders).Error
	return orders, err
}
