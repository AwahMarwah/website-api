package category

import (
	"website-api/library/helper/filter"
	"website-api/model/category"
)

func (r *repo) GetCategory(reqQuery *category.FilterCategory) (resData []*category.ListCategoryResponse, count int64, err error) {
	scope := filter.FilterCategorySearch(reqQuery.Search)

	if err = r.db.Model(&category.Category{}).
		Scopes(scope).
		Count(&count).Error; err != nil {
		return resData, count, err
	}

	err = r.db.Model(&category.Category{}).
		Scopes(scope).
		Limit(reqQuery.Limit).
		Offset(reqQuery.Offset).
		Find(&resData).Error
	return resData, count, err
}