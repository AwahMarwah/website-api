package cart

import (
	cartRepo "website-api/repository/cart"
	"website-api/service/cart"

	"gorm.io/gorm"
)

type controller struct {
	cartService cart.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{cartService: cart.NewService(cartRepo.NewRepo(db))}
}
