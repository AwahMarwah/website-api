package user

import (
	userModel "website-api/model/user"
	"website-api/repository/role"
	"website-api/repository/user"
	userAddressRepo "website-api/repository/user_address"
)

type (
	IService interface {
		List(reqQuery *userModel.ListUserReqQuery) (resData []userModel.ListUserResponse, count int64, err error)
		Detail(reqPath *userModel.ReqPath) (resData userModel.DetailResponse, statusCode int, err error)
		SignIn(reqBody *userModel.SignInRequest) (resData userModel.SignInResponse, statusCode int, err error)
		SignUp(reqBody *userModel.RegisterRequest) (statusCode int, err error)
		VerifyEmail(reqBody userModel.VerifyEmailRequest) (err error)
		Update(req *userModel.UserUpdateRequest) (statusCode int, err error)
		Seed() (err error)
		SignOut(userId string) (err error)
	}

	service struct {
		userRepo        user.IRepo
		roleRepo        role.IRepo
		userAddressRepo userAddressRepo.IRepo
	}
)

func NewService(userRepo user.IRepo, roleRepo role.IRepo, userAddressRepo userAddressRepo.IRepo) IService {
	return &service{
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		userAddressRepo: userAddressRepo,
	}
}
