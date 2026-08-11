package role

import (
	modelRole "website-api/model/role"
)

func (r *repo) Take(selectParams []string, conditions *modelRole.Role) (role modelRole.Role, err error) {
	return role, r.db.Select(selectParams).Take(&role, conditions).Error
}
