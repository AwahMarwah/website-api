package master

import "website-api/model/master"

func (s *service) GetDistrict(reqQuery *master.GetListDistrictsRequest) (resData []master.District, count int64, err error) {
	resData, count, err = s.masterRepo.FindDistrict(reqQuery)
	if err != nil {
		return
	}
	return
}
