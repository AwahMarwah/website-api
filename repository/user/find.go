package user

import (
	"website-api/library/helper/filter"
	userModel "website-api/model/user"
)

func (r *repo) Find(reqQuery *userModel.ListUserReqQuery) (resData []userModel.ListUserResponse, count int64, err error) {
	resData = make([]userModel.ListUserResponse, 0)
	q := r.db.Model(&userModel.User{}).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Select("users.*, roles.name as role_name").
		Scopes(filter.FilterUserSearch(reqQuery.Search)).
		Count(&count)
	if err = q.Limit(reqQuery.Limit).Offset(reqQuery.Offset).Scan(&resData).Error; err != nil {
		return
	}
	return
}