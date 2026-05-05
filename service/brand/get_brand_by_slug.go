package brand

import brandModel "website-api/model/brand"

func (s *service) GetBrandBySlug(reqBody *brandModel.FilterBrandReq) (brand brandModel.Brand, err error) {
	return s.brandRepo.GetBrandBySlug(reqBody)
}
