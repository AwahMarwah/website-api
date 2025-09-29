package user

import userModel "website-api/model/user"

func (r *repo) Create(reqBody *userModel.User) error {
	return r.db.Create(reqBody).Error
}
