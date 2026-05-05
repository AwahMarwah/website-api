package product

import (
	"fmt"
	"time"
	libStruct "website-api/library"
	lib "website-api/library/cache"
	productModel "website-api/model/product"
)

func (s *service) GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error) {
	cacheKey := lib.GenerateCacheKey(libStruct.GetStructName(productModel.Product{}), reqQuery)

	// Check Cache
	err = s.cache.Get(cacheKey, &resData)
	if err == nil {
		fmt.Printf("CACHE HIT | key=%s\n", cacheKey)
		return resData, count, nil
	}

	fmt.Printf("CACHE MISS | key=%s\n", cacheKey)

	// Get from DB
	resData, count, err = s.productRepo.GetProduct(reqQuery)
	if err != nil {
		return nil, count, err
	}

	// Save to Redis
	_ = s.cache.Set(cacheKey, resData, 5*time.Minute)
	fmt.Printf("CACHE SET | key=%s | ttl=5m\n", cacheKey)

	return resData, count, nil
}
