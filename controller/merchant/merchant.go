package merchant

import (
	merchantRepo "website-api/repository/merchant"
	"website-api/service/merchant"

	"gorm.io/gorm"
)

type controller struct {
	merchantService merchant.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{merchantService: merchant.NewService(merchantRepo.NewRepo(db))}
}