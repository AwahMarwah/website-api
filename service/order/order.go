package order

import (
	"website-api/database/transaction"
	"website-api/model/order"
	"website-api/repository/product"
	product_variant "website-api/repository/product-variant"
)

type (
	IService interface {
		Checkout(req *order.CheckoutRequest) error
	}

	service struct {
		productRepo        product.IRepo
		productVariantRepo product_variant.IRepo
		txManager          transaction.ITransactionManager
	}
)

func NewService(productRepo product.IRepo, productVariantRepo product_variant.IRepo, txManager transaction.ITransactionManager) IService {
	return &service{
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		txManager:          txManager,
	}
}
