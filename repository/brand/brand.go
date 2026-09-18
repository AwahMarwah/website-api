package brand

import (
	brandModel "website-api/model/brand"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		GetBrand(reqQuery *brandModel.BrandReqQuery) (resData []brandModel.ListBrandResponse, count int64, err error)
		GetBrandBySlug(reqBody *brandModel.FilterBrandReq) (brand brandModel.Brand, err error)
		Create(brand brandModel.Brand) error
		FindByID(id string) (brand brandModel.Brand, err error)
		Update(id string, values map[string]any) error
		Delete(id string) error
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
