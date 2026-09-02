package order

import (
	"website-api/database/transaction"
	"website-api/model/order"
	"website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	orderRepo "website-api/repository/order"
	userRepo "website-api/repository/user"
	"website-api/third-party/provider/midtrans"

	"github.com/hibiken/asynq"
)

type (
	QueueClient interface {
		Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
	}

	IService interface {
		Checkout(req *order.CheckoutRequest) (order.CheckoutResponse, int, error)
		CreatePaymentLink(orderID string) (order.PaymentLinkResponse, int, error)
		HandleNotification(payload midtrans.NotificationPayload) (order.NotificationResponse, int, error)
		CancelExpiredOrders() (int, error)
	}

	service struct {
		productRepo        product.IRepo
		productVariantRepo product_variant.IRepo
		orderRepo          orderRepo.IRepo
		userRepo           userRepo.IRepo
		txManager          transaction.ITransactionManager
		midtransProvider   midtrans.Provider
		queueClient        QueueClient
	}
)

func NewService(
	productRepo product.IRepo,
	productVariantRepo product_variant.IRepo,
	orderRepo orderRepo.IRepo,
	userRepo userRepo.IRepo,
	txManager transaction.ITransactionManager,
	midtransProvider midtrans.Provider,
	queueClient QueueClient,
) IService {
	return &service{
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		orderRepo:          orderRepo,
		userRepo:           userRepo,
		txManager:          txManager,
		midtransProvider:   midtransProvider,
		queueClient:        queueClient,
	}
}
