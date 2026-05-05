package category

import (
	categoryRepo "website-api/repository/category"
	"website-api/service/category"

	"gorm.io/gorm"
)

type controller struct {
	categoryService category.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{categoryService: category.NewService(categoryRepo.NewRepo(db))}
}
