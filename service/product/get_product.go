package product

import (
	"time"
	libStruct "website-api/library"
	lib "website-api/library/cache"
	productModel "website-api/model/product"
)

func (s *service) GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error) {
	cacheKey := lib.GenerateCacheKey(libStruct.GetStructName(productModel.Product{}), reqQuery)

	// Check Cache
	var cached productModel.ProductListCache
	if err = s.cache.Get(cacheKey, &cached); err == nil {
		return cached.Data, cached.Count, nil
	}

	// Get from DB
	resData, count, err = s.productRepo.GetProduct(reqQuery)
	if err != nil {
		return nil, count, err
	}

	// Save to Redis
	_ = s.cache.Set(cacheKey, productModel.ProductListCache{Data: resData, Count: count}, 5*time.Minute)

	return resData, count, nil
}