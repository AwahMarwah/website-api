package user_address

import userAddressModel "website-api/model/user_address"

func (r *repo) Update(id *string, values *map[string]any) (err error) {
	return r.db.Model(userAddressModel.UserAddress{ID: *id}).Updates(values).Error
}
