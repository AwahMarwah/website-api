package user_address

import "website-api/model/user_address"

func (r *repo) Find(userID string) (resData []user_address.UserAddress, err error) {
	resData = make([]user_address.UserAddress, 0)
	return resData, r.db.Model(&user_address.UserAddress{}).Where("user_id = ?", userID).Order("created_at desc").Limit(4).Find(&resData).Error
}
