package merchant

import (
	"website-api/database/transaction"
	merchantModel "website-api/model/merchant"
	"website-api/repository/merchant"
	"website-api/repository/order"
	"website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	"website-api/repository/settlement"
)

type (
	IService interface {
		List() ([]merchantModel.MerchantResponse, int, error)
		Detail(id string) (merchantModel.MerchantResponse, int, error)
		Register(req *merchantModel.CreateMerchantReq, userID string) (int, error)
		GetMy(userID string) (merchantModel.MerchantResponse, int, error)
		UpdateMy(userID string, req *merchantModel.UpdateMerchantReq) (int, error)
		Approve(id string) (int, error)

		// Seller panel. Semua method menurunkan merchant dari userID, bukan dari request body.
		ListProducts(userID string, reqQuery *merchantModel.ListSellerProductReqQuery) (resData []merchantModel.SellerProductResponse, count int64, statusCode int, err error)
		GetProduct(userID, productID string) (merchantModel.SellerProductResponse, int, error)
		CreateProduct(userID string, req *merchantModel.CreateSellerProductReq) (productID string, statusCode int, err error)
		UpdateProduct(userID, productID string, req *merchantModel.UpdateSellerProductReq) (int, error)
		DeleteProduct(userID, productID string) (int, error)

		ListVariants(userID, productID string) ([]merchantModel.SellerVariantResponse, int, error)
		CreateVariant(userID, productID string, req *merchantModel.CreateSellerVariantReq) (variantID string, statusCode int, err error)
		UpdateVariant(userID, variantID string, req *merchantModel.UpdateSellerVariantReq) (int, error)
		DeleteVariant(userID, variantID string) (int, error)

		ListOrders(userID string, reqQuery *merchantModel.ListSellerOrderReqQuery) (resData []merchantModel.SellerOrderResponse, count int64, statusCode int, err error)
		GetOrder(userID, orderID string) (merchantModel.SellerOrderResponse, int, error)
		Stats(userID string) (merchantModel.SellerStatsResponse, int, error)
	}

	service struct {
		merchantRepo       merchant.IRepo
		productRepo        product.IRepo
		productVariantRepo product_variant.IRepo
		orderRepo          order.IRepo
		settlementRepo     settlement.IRepo
		txManager          transaction.ITransactionManager
	}
)

func NewService(
	merchantRepo merchant.IRepo,
	productRepo product.IRepo,
	productVariantRepo product_variant.IRepo,
	orderRepo order.IRepo,
	settlementRepo settlement.IRepo,
	txManager transaction.ITransactionManager,
) IService {
	return &service{
		merchantRepo:       merchantRepo,
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		orderRepo:          orderRepo,
		settlementRepo:     settlementRepo,
		txManager:          txManager,
	}
}
