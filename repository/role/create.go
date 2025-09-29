package role

import modelRole "website-api/model/role"

func (r *repo) Create(reqBody *modelRole.Role) error {
	return r.db.Create(reqBody).Error
}
