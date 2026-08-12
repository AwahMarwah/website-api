package user_address

import userAddressModel "website-api/model/user_address"

func (r *repo) Create(reqBody *userAddressModel.UserAddress) error {
	return r.db.Create(reqBody).Error
}
