package product

import (
	"website-api/cache"
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	"website-api/service/product"
	"website-api/service/upload"
	minioProvider "website-api/third-party/provider/minio"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type controller struct {
	productService product.IService
	uploadService  upload.IService
}

func NewController(db *gorm.DB, redis *redis.Client, minio minioProvider.Provider) *controller {
	redisCache := cache.NewRedisCache(redis)
	return &controller{
		productService: product.NewService(productRepo.NewRepo(db), product_variant.NewRepo(db), redisCache),
		uploadService:  upload.NewService(minio, nil, productRepo.NewRepo(db)),
	}
}
