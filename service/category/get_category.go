package category

import "website-api/model/category"

func (s *service) GetCategory(reqQuery *category.FilterCategory) (resData []*category.ListCategoryResponse, count int64, err error) {
	resData, count, err = s.categoryRepo.GetCategory(reqQuery)
	if err != nil {
		return resData, count, err
	}
	return
}
