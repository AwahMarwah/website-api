package product

import (
	"website-api/cache"
	"website-api/database/transaction"
	productModel "website-api/model/product"
	"website-api/repository/product"
	product_variant "website-api/repository/product-variant"
)

type (
	IService interface {
		GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error)
		GetProductDetail(id string) (resData productModel.ProductDetailResponse, err error)
		CreateProduct(req *productModel.CreateProductReq) (int, error)
		UpdateProduct(id string, req *productModel.UpdateProductReq) (int, error)
		DeleteProduct(id string) (int, error)
		SetPrimaryImage(productID, imageID string) (int, error)
		DeleteImage(productID, imageID string) (int, error)
		ReorderImages(productID string, orders []productModel.ImageSortOrder) (int, error)
	}

	service struct {
		productRepo        product.IRepo
		productVariantRepo product_variant.IRepo
		cache              cache.Cache
		txManager          transaction.ITransactionManager
	}
)

func NewService(productRepo product.IRepo, productVariantRepo product_variant.IRepo, redis cache.Cache, txManager transaction.ITransactionManager) IService {
	return &service{
		productRepo:        productRepo,
		productVariantRepo: productVariantRepo,
		cache:              redis,
		txManager:          txManager,
	}
}
