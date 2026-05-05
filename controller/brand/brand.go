package brand

import (
	brandRepo "website-api/repository/brand"
	"website-api/service/brand"

	"gorm.io/gorm"
)

type controller struct {
	brandService brand.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{brandService: brand.NewService(brandRepo.NewRepo(db))}
}
