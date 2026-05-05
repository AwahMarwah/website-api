package brand

import (
	brandModel "website-api/model/brand"
)

func (r *repo) GetBrandBySlug(reqBody *brandModel.FilterBrandReq) (brand brandModel.Brand, err error) {
	return brand, r.db.Debug().Model(&brandModel.Brand{}).Where("slug = ?", reqBody.Slug).First(&brand).Error
}
