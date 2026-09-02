package order

import (
	"website-api/database/transaction"
	orderRepo "website-api/repository/order"
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	userRepo "website-api/repository/user"
	"website-api/service/order"
	midtransProvider "website-api/third-party/provider/midtrans"
	"website-api/worker"

	"gorm.io/gorm"
)

type controller struct {
	orderService order.IService
	midtrans     *midtransProvider.Client
}

func NewController(db *gorm.DB) *controller {
	midtransClient := midtransProvider.NewClient()
	return &controller{
		orderService: order.NewService(
			productRepo.NewRepo(db),
			product_variant.NewRepo(db),
			orderRepo.NewRepo(db),
			userRepo.NewRepo(db),
			transaction.NewTransactionManager(db),
			midtransClient,
			worker.NewRedisClient(),
		),
		midtrans: midtransClient,
	}
}
