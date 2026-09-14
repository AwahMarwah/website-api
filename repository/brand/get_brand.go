package brand

import brandModel "website-api/model/brand"

func (r *repo) GetBrand(reqQuery *brandModel.BrandReqQuery) (resData []brandModel.ListBrandResponse, count int64, err error) {
	resData = make([]brandModel.ListBrandResponse, 0)
	q := r.db.Model(&brandModel.Brand{})
	if err := q.Count(&count).Error; err != nil {
		return nil, count, err
	}
	err = q.Order("name ASC").Limit(reqQuery.Limit).Offset(reqQuery.Offset).Find(&resData).Error
	return resData, count, err
}