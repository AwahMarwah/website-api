package review

import (
	orderRepo "website-api/repository/order"
	reviewRepo "website-api/repository/review"
	reviewService "website-api/service/review"

	"gorm.io/gorm"
)

type controller struct {
	reviewService reviewService.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{
		reviewService: reviewService.NewService(reviewRepo.NewRepo(db), orderRepo.NewRepo(db)),
	}
}