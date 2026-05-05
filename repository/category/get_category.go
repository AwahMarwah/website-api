package category

import (
	"website-api/library/helper/filter"
	"website-api/model/category"
)

func (r *repo) GetCategory(reqQuery *category.FilterCategory) (resData []*category.ListCategoryResponse, count int64, err error) {
	return resData, count, r.db.Debug().Model(&category.Category{}).Scopes(filter.FilterCategory(reqQuery.Search)).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Find(&resData).Error
}
