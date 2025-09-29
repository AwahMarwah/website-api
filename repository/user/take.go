package user

import modelUser "website-api/model/user"

func (r *repo) Take(selectParams []string, conditions *modelUser.User) (user modelUser.User, err error) {
	return user, r.db.Debug().Select(selectParams).Take(&user, conditions).Error
}
