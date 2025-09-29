package auth

import (
	authModel "website-api/model/auth"
)

func (r *repo) Update(id *string, values *map[string]any) (err error) {
	return r.db.Debug().Model(authModel.PasswordResetToken{Id: *id}).Updates(values).Error
}
