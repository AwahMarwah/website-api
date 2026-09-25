package product

import (
	"time"

	"website-api/model/brand"
	"website-api/model/merchant"
	"website-api/model/review"
)

type Product struct {
	Id          string
	BrandId     string
	MerchantId  string
	Sku         string
	Name        string
	Slug        string
	Description string
	BasePrice   float32
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Brand       *brand.Brand              `gorm:"foreignKey:BrandId"`
	Merchant    *merchant.Merchant        `gorm:"foreignKey:MerchantId"`
	Reviews     []review.Review           `gorm:"foreignKey:ProductID"`
	Images      []ProductImage            `gorm:"foreignKey:ProductID"`
}
