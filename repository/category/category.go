package category

import (
	"website-api/model/category"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		GetCategory(reqQuery *category.FilterCategory) (resData []*category.ListCategoryResponse, count int64, err error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db}
}
