package product_variant

import (
	"time"

	"website-api/model/product"
)

type ProductVariant struct {
	ID          string
	ProductID   string
	Sku         string
	VariantName string
	Price       float64
	Stock       int
	Weight      float32
	IsActive    bool
	CreatedAt   time.Time
	Product     *product.Product `gorm:"foreignKey:ProductID"`
}
