package user_address

import (
	"errors"
	userAddressModel "website-api/model/user_address"
)

func (s *service) UpdateByID(req *userAddressModel.ReqUpdateUserAddress) error {
	// check existing data
	userAddress, err := s.userAddressRepo.Take([]string{"id"}, &userAddressModel.UserAddress{ID: req.Path.ID})
	if err != nil {
		return err
	}

	if userAddress.ID == "" {
		return errors.New("record not found")
	}

	values := map[string]any{
		"recipient_name": req.Body.RecipientName,
		"phone_number":   req.Body.PhoneNumber,
		"full_address":   req.Body.FullAddress,
		"province_id":    req.Body.ProvinceID,
		"city_id":        req.Body.CityID,
		"district_id":    req.Body.DistrictID,
		"subdistrict_id": req.Body.SubdistrictID,
		"postal_code":    req.Body.PostalCode,
	}

	return s.userAddressRepo.Update(&req.Path.ID, &values)
}
