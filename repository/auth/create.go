package auth

import authModel "website-api/model/auth"

func (r *repo) Create(auth *authModel.PasswordResetToken) (err error) {
	return r.db.Create(auth).Error
}
