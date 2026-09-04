package master

import (
	"website-api/library/helper/filter"
	city_shipping "website-api/model/city-shipping"
	"website-api/model/master"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		TakeProvince(selectParams []string, conditions *master.Province) (province master.Province, err error)
		TakeCity(selectParams []string, conditions *master.City) (city master.City, err error)
		TakeDistrict(selectParams []string, conditions *master.District) (district master.District, err error)
		TakeSubdistrict(selectParams []string, conditions *master.Subdistrict) (subdistrict master.Subdistrict, err error)
		TakeCityShippingMapping(selectParams []string, conditions *city_shipping.CityShippingMapping) (city_shipping city_shipping.CityShippingMapping, err error)
		CreateCityShippingMapping(reqBody *city_shipping.CityShippingMapping) error
		FindProvince(reqQuery *master.GetListProvinceRequest) (resData []master.ListProvinceResponse, count int64, err error)
		FindCities(reqQuery *master.GetListCitiesRequest) (resData []master.City, count int64, err error)
		FindDistrict(reqQuery *master.GetListDistrictsRequest) (resData []master.ListDistrictResponse, count int64, err error)
		FindSubdistrict(reqQuery *master.GetListSubdistrictsRequest) (resData []master.ListSubdistrictResponse, count int64, err error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) TakeProvince(selectParams []string, conditions *master.Province) (province master.Province, err error) {
	return province, r.db.Select(selectParams).Take(&province, conditions).Error
}

func (r *repo) TakeCity(selectParams []string, conditions *master.City) (city master.City, err error) {
	return city, r.db.Select(selectParams).Take(&city, conditions).Error
}

func (r *repo) TakeDistrict(selectParams []string, conditions *master.District) (district master.District, err error) {
	return district, r.db.Select(selectParams).Take(&district, conditions).Error
}

func (r *repo) TakeSubdistrict(selectParams []string, conditions *master.Subdistrict) (subdistrict master.Subdistrict, err error) {
	return subdistrict, r.db.Select(selectParams).Take(&subdistrict, conditions).Error
}

func (r *repo) TakeCityShippingMapping(selectParams []string, conditions *city_shipping.CityShippingMapping) (city_shipping city_shipping.CityShippingMapping, err error) {
	return city_shipping, r.db.Select(selectParams).Take(&city_shipping, conditions).Error
}

func (r *repo) CreateCityShippingMapping(reqBody *city_shipping.CityShippingMapping) error {
	return r.db.Create(reqBody).Error
}

func (r *repo) FindProvince(reqQuery *master.GetListProvinceRequest) (resData []master.ListProvinceResponse, count int64, err error) {
	resData = make([]master.ListProvinceResponse, 0)

	return resData, count, r.db.Model(&master.Province{}).Scopes(filter.FilterProvinceSearch(reqQuery.Search)).Count(&count).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Order("id").Find(&resData).Error
}

func (r *repo) FindCities(reqQuery *master.GetListCitiesRequest) (resData []master.City, count int64, err error) {
	resData = make([]master.City, 0)

	return resData, count, r.db.Debug().Model(&master.City{}).Scopes(filter.FilterCitiesSearch(reqQuery.Search)).Count(&count).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Order("id").Find(&resData).Error
}

func (r *repo) FindDistrict(reqQuery *master.GetListDistrictsRequest) (resData []master.ListDistrictResponse, count int64, err error) {
	resData = make([]master.ListDistrictResponse, 0)

	return resData, count, r.db.Debug().Model(&master.District{}).Scopes(filter.FilterDistrictSearch(reqQuery.Search)).Count(&count).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Order("id").Find(&resData).Error
}

func (r *repo) FindSubdistrict(reqQuery *master.GetListSubdistrictsRequest) (resData []master.ListSubdistrictResponse, count int64, err error) {
	resData = make([]master.ListSubdistrictResponse, 0)

	return resData, count, r.db.Debug().Model(&master.Subdistrict{}).Scopes(filter.FilterSubdistrictSearch(reqQuery.Search)).Count(&count).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Order("id").Find(&resData).Error
}
