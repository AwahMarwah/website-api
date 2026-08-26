package master

import (
	masterModel "website-api/model/master"
)

func (s *service) GetProvince(reqQuery *masterModel.GetListProvinceRequest) (resData []masterModel.ListProvinceResponse, count int64, err error) {
	resData, count, err = s.masterRepo.FindProvince(reqQuery)
	if err != nil {
		return nil, count, err
	}
	return resData, count, err
}
