package brand

import brandModel "website-api/model/brand"

func (r *repo) GetBrand(reqQuery *brandModel.BrandReqQuery) (resData []brandModel.ListBrandResponse, count int64, err error) {
	return resData, count, r.db.Debug().Model(&brandModel.Brand{}).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Find(&resData).Error
}
