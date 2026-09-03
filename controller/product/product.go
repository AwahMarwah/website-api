package product

import (
	"website-api/cache"
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	"website-api/service/product"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type controller struct {
	productService product.IService
}

func NewController(db *gorm.DB, redis *redis.Client) *controller {
	redisCache := cache.NewRedisCache(redis)
	return &controller{productService: product.NewService(productRepo.NewRepo(db), product_variant.NewRepo(db), redisCache)}
}
