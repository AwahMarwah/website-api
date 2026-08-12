package user_address

import userAddressModel "website-api/model/user_address"

func (r *repo) Delete(id string) error {
	return r.db.Delete(&userAddressModel.UserAddress{}, "id = ?", id).Error
}
