package role

import roleModel "website-api/model/role"

func (r *repo) Update(id *string, values *map[string]any) (err error) {
	return r.db.Model(roleModel.Role{Id: *id}).Updates(values).Error
}
