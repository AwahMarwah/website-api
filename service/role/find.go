package role

import modelRole "website-api/model/role"

func (s *service) Find() (roles []modelRole.Role, err error) {
	return s.roleRepo.Find()
}
