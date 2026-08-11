package user_address

import "website-api/model/user_address"

func (r *repo) Take(selectParams []string, conditions *user_address.UserAddress) (address user_address.UserAddress, err error) {
	return address, r.db.Select(selectParams).Take(&address, conditions).Error
}
