package order

import (
	"website-api/database/transaction"
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	"website-api/service/order"

	"gorm.io/gorm"
)

type controller struct {
	orderService order.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{order.NewService(productRepo.NewRepo(db), product_variant.NewRepo(db), transaction.NewTransactionManager(db))}
}
