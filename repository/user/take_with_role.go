package user

import modelUser "website-api/model/user"

// TakeWithRole mengambil satu user beserta relasi role-nya dalam satu query Preload.
func (r *repo) TakeWithRole(selectParams []string, conditions *modelUser.User) (user modelUser.User, err error) {
	return user, r.db.
		Select(selectParams).
		Preload("Role").
		Take(&user, conditions).Error
}