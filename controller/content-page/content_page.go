package content_page

import (
	"website-api/cache"
	contentPageRepo "website-api/repository/content-page"
	"website-api/service/content_page"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type controller struct {
	contentPageService content_page.IService
}

func NewController(db *gorm.DB, redis *redis.Client) *controller {
	redisCache := cache.NewRedisCache(redis)
	return &controller{contentPageService: content_page.NewService(contentPageRepo.NewRepo(db), redisCache)}
}
