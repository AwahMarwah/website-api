package user

import userModel "website-api/model/user"

func (r *repo) Find(reqQuery *userModel.ListUserReqQuery) (resData []userModel.ListUserResponse, count int64, err error) {
	resData = make([]userModel.ListUserResponse, 0)
	return resData, count, r.db.Model(&userModel.User{}).Count(&count).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Scan(&resData).Error
}
