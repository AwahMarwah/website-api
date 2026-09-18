package brand

import (
	brandModel "website-api/model/brand"
	"website-api/repository/brand"
)

type (
	IService interface {
		GetBrand(reqQuery *brandModel.BrandReqQuery) (resData []brandModel.ListBrandResponse, count int64, err error)
		GetBrandBySlug(reqBody *brandModel.FilterBrandReq) (brand brandModel.Brand, err error)
		Create(req *brandModel.CreateBrandReq) (int, error)
		Update(id string, req *brandModel.UpdateBrandReq) (int, error)
		Delete(id string) (int, error)
	}

	service struct {
		brandRepo brand.IRepo
	}
)

func NewService(brandRepo brand.IRepo) IService {
	return &service{brandRepo: brandRepo}
}
