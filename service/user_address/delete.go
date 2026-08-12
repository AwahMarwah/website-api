package user_address

import userAddressModel "website-api/model/user_address"

func (s *service) Delete(reqPath userAddressModel.ReqPath) error {
	return s.userAddressRepo.Delete(reqPath.ID)
}
