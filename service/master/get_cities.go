package master

import "website-api/model/master"

func (s *service) GetCities(reqQuery *master.GetListCitiesRequest) (resData []master.City, count int64, err error) {
	resData, count, err = s.masterRepo.FindCities(reqQuery)
	if err != nil {
		return
	}
	return resData, count, err
}
