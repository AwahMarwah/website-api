package product_variant

import (
	product_variant "website-api/model/product-variant"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		FindByID(id string) (productVariant product_variant.ProductVariant, err error)
		Update(productVariant product_variant.ProductVariant) error
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) FindByID(id string) (productVariant product_variant.ProductVariant, err error) {
	return productVariant, r.db.Where("id = ?", id).First(&productVariant).Error
}

func (r *repo) Update(productVariant product_variant.ProductVariant) error {
	return r.db.Model(&product_variant.ProductVariant{}).Where("id = ?", productVariant.ID).Update("stock", productVariant.Stock).Error
}

func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{db: tx}
}
