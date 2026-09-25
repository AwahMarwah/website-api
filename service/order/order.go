package order

import (
	"website-api/database/transaction"
	"website-api/model/order"
	merchantRepo "website-api/repository/merchant"
	"website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	orderRepo "website-api/repository/order"
	reviewRepo "website-api/repository/review"
	userRepo "website-api/repository/user"
	userAddressRepo "website-api/repository/user_address"
	"website-api/third-party/provider/midtrans"
	"website-api/third-party/provider/rajaongkir"

	"github.com/hibiken/asynq"
)

type (
	QueueClient interface {
		Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
	}

	IService interface {
		Checkout(req *order.CheckoutRequest) (order.CheckoutResponse, int, error)
		CreatePaymentLink(orderID, userID, roleName string) (order.PaymentLinkResponse, int, error)
		HandleNotification(payload midtrans.NotificationPayload) (order.NotificationResponse, int, error)
		CancelExpiredOrders() (int, error)
		List(req *order.ListOrderReqQuery) (resData []order.OrderResponse, count int64, statusCode int, err error)
		Detail(id, userID, roleName string) (resData order.OrderResponse, statusCode int, err error)
		ListAdmin(req *order.ListOrderReqQuery) (resData []order.OrderResponse, count int64, statusCode int, err error)
		UpdateStatusAdmin(id, status string) (int, error)
		CancelOrder(orderID, userID string) (int, error)
	}

	service struct {
		productRepo        product.IRepo
		productVariantRepo product_variant.IRepo
		orderRepo          orderRepo.IRepo
		userRepo           userRepo.IRepo
		userAddressRepo    userAddressRepo.IRepo
		merchantRepo       merchantRepo.IRepo
		reviewRepo         reviewRepo.IRepo
		txManager          transaction.ITransactionManager
		midtransProvider   midtrans.Provider
		rajaOngkir         rajaongkir.Provider
		queueClient        QueueClient
	}
)

func NewService(
	productRepo product.IRepo,
	productVariantRepo product_variant.IRepo,
	orderRepo orderRepo.IRepo,
	userRepo userRepo.IRepo,
	userAddressRepo userAddressRepo.IRepo,
	merchantRepo merchantRepo.IRepo,
	reviewRepo reviewRepo.IRepo,
	txManager transaction.ITransactionManager,
	midtransProvider midtrans.Provider,
	rajaOngkir rajaongkir.Provider,
	queueClient QueueClient,
) IService {
	return &service{
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		orderRepo:          orderRepo,
		userRepo:           userRepo,
		userAddressRepo:    userAddressRepo,
		merchantRepo:       merchantRepo,
		reviewRepo:         reviewRepo,
		txManager:          txManager,
		midtransProvider:   midtransProvider,
		rajaOngkir:         rajaOngkir,
		queueClient:        queueClient,
	}
}
