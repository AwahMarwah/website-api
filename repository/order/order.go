package order

import (
	"website-api/model/order"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(order order.OrderItem) error
	}

	repo struct {
		db gorm.DB
	}
)

func NewRepo(db gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) Create(order order.OrderItem) error {
	return r.db.Create(order).Error
}
