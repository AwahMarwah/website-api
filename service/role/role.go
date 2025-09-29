package role

import (
	modelRole "website-api/model/role"
	roleRepo "website-api/repository/role"
)

type (
	IService interface {
		Create(reqBody *modelRole.RoleCreateReqBody) (statusCode int, err error)
		Detail(reqPath *modelRole.ReqPath) (resData modelRole.DetailResponse, statusCode int, err error)
		Find() (roles []modelRole.Role, err error)
		Update(req *modelRole.UpdateReq) (statusCode int, err error)
		Delete(reqPath *modelRole.ReqPath) (statusCode int, err error)
		Seed() (err error)
	}

	service struct {
		roleRepo roleRepo.IRepo
	}
)

func NewService(roleRepo roleRepo.IRepo) IService {
	return &service{roleRepo: roleRepo}
}
