package auth

import authModel "website-api/model/auth"

func (r *repo) Delete(conditions map[string]interface{}) error {
	return r.db.Where(conditions).Delete(&authModel.PasswordResetToken{}).Error
}
