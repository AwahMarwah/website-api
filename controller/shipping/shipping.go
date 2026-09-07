package shipping

import (
	merchantRepo "website-api/repository/merchant"
	productRepo "website-api/repository/product"
	product_variant "website-api/repository/product-variant"
	userAddressRepo "website-api/repository/user_address"
	"website-api/service/shipping"
	"website-api/third-party/provider/rajaongkir"

	"gorm.io/gorm"
)

type controller struct {
	shippingService shipping.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{
		shippingService: shipping.NewService(
			userAddressRepo.NewRepo(db),
			product_variant.NewRepo(db),
			productRepo.NewRepo(db),
			merchantRepo.NewRepo(db),
			rajaongkir.NewClient(),
		),
	}
}