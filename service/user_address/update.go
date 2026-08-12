package user_address

import (
	userAddressModel "website-api/model/user_address"
)

func (s *service) Update(req *userAddressModel.UpdateUserAddress) error {
	values := map[string]any{
		"is_primary": req.Body.IsPrimary,
	}
	return s.userAddressRepo.Update(&req.Path.ID, &values)
}
