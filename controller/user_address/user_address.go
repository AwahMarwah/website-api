package user_address

import (
	"website-api/cache"
	"website-api/repository/master"
	userAddressRepo "website-api/repository/user_address"
	"website-api/service/user_address"
	"website-api/third-party/provider/rajaongkir"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type controller struct {
	userAddressService user_address.IService
}

func NewController(db *gorm.DB, redis *redis.Client) *controller {
	redisCache := cache.NewRedisCache(redis)
	rajaOngkir := rajaongkir.NewClient()
	return &controller{userAddressService: user_address.NewService(userAddressRepo.NewRepo(db), master.NewRepo(db), redisCache, rajaOngkir)}
}
