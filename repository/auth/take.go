package auth

import authModel "website-api/model/auth"

func (r *repo) Take(selectParams []string, condition *authModel.PasswordResetToken) (result authModel.PasswordResetToken, err error) {
	return result, r.db.Select(selectParams).Where(condition).Take(&result).Error
}
