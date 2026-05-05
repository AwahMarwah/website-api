package category

import (
	"website-api/model/category"
	categoryRepo "website-api/repository/category"
)

type (
	IService interface {
		GetCategory(reqQuery *category.FilterCategory) (resData []*category.ListCategoryResponse, count int64, err error)
	}

	service struct {
		categoryRepo categoryRepo.IRepo
	}
)

func NewService(categoryRepo categoryRepo.IRepo) IService {
	return &service{categoryRepo: categoryRepo}
}
