package user_address

import (
	"website-api/cache"
	userAddressRepo "website-api/repository/user_address"
	"website-api/service/user_address"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type controller struct {
	userAddressService user_address.IService
}

func NewController(db *gorm.DB, redis *redis.Client) *controller {
	redisCache := cache.NewRedisCache(redis)
	return &controller{userAddressService: user_address.NewService(userAddressRepo.NewRepo(db), redisCache)}
}
