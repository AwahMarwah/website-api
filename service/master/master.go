package master

import (
	masterModel "website-api/model/master"
	"website-api/repository/master"
)

type (
	IService interface {
		GetProvince(reqQuery *masterModel.GetListProvinceRequest) (resData []masterModel.ListProvinceResponse, count int64, err error)
		GetCities(reqQuery *masterModel.GetListCitiesRequest) (resData []masterModel.City, count int64, err error)
		GetDistrict(reqQuery *masterModel.GetListDistrictsRequest) (resData []masterModel.ListDistrictResponse, count int64, err error)
		GetSubdistrict(reqQuery *masterModel.GetListSubdistrictsRequest) (resData []masterModel.ListSubdistrictResponse, count int64, err error)
	}

	service struct {
		masterRepo master.IRepo
	}
)

func NewService(masterRepo master.IRepo) IService {
	return &service{
		masterRepo: masterRepo,
	}
}
