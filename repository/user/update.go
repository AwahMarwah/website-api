package user

import userModel "website-api/model/user"

func (r *repo) Update(id *string, values *map[string]any) (err error) {
	return r.db.Debug().Model(userModel.User{Id: *id}).Updates(values).Error
}
