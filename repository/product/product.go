package product

import (
	productModel "website-api/model/product"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		FindByID(id string) (resData productModel.Product, err error)
		GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error)
		Update(product productModel.Product) error
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
