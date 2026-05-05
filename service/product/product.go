package product

import (
	"website-api/cache"
	productModel "website-api/model/product"
	"website-api/repository/product"
)

type (
	IService interface {
		GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error)
	}

	service struct {
		productRepo product.IRepo
		cache       cache.Cache
	}
)

func NewService(productRepo product.IRepo, redis cache.Cache) IService {
	return &service{
		productRepo: productRepo,
		cache:       redis,
	}
}
