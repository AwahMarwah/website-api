package user_address

import (
	"website-api/cache"
	userAddressModel "website-api/model/user_address"
	userAddressRepo "website-api/repository/user_address"
)

type (
	IService interface {
		Find(userID string) (resData []userAddressModel.UserAddress, err error)
		Create(reqBody *userAddressModel.CreateUserAddressRequest) error
		Detail(reqPath *userAddressModel.ReqPath) (resData userAddressModel.UserAddress, err error)
		Delete(reqPath userAddressModel.ReqPath) error
		Update(reqPath *userAddressModel.UpdateUserAddress) error
	}

	service struct {
		userAddressRepo userAddressRepo.IRepo
		redis           cache.Cache
	}
)

func NewService(userAddressRepo userAddressRepo.IRepo, redis cache.Cache) IService {
	return &service{
		userAddressRepo: userAddressRepo,
		redis:           redis,
	}
}
