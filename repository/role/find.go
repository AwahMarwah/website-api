package role

import modelRole "website-api/model/role"

func (r *repo) Find() (roles []modelRole.Role, err error) {
	return roles, r.db.Where("deleted_at IS NULL").Find(&roles).Error
}
