package brand

import brandModel "website-api/model/brand"

func (s *service) GetBrand(reqQuery *brandModel.BrandReqQuery) (resData []brandModel.ListBrandResponse, count int64, err error) {
	resData, count, err = s.brandRepo.GetBrand(reqQuery)
	if err != nil {
		return resData, count, err
	}
	return
}
