package user_address

import userAddressModel "website-api/model/user_address"

func (s *service) Detail(reqPath *userAddressModel.ReqPath) (resData userAddressModel.UserAddress, err error) {
	userAddress, err := s.userAddressRepo.Take([]string{"id", "user_id", "recipient_name", "phone_number", "full_address", "city", "postal_code", "is_primary"}, &userAddressModel.UserAddress{ID: reqPath.ID})
	if err != nil {
		return resData, err
	}
	resData = userAddress
	return resData, nil
}
