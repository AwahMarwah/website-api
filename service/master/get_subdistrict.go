package master

import "website-api/model/master"

func (s *service) GetSubdistrict(reqQuery *master.GetListSubdistrictsRequest) (resData []master.ListSubdistrictResponse, count int64, err error) {
	resData, count, err = s.masterRepo.FindSubdistrict(reqQuery)
	if err != nil {
		return
	}
	return
}
