package user_address

import (
	"website-api/cache"
	userAddressModel "website-api/model/user_address"
	"website-api/repository/master"
	userAddressRepo "website-api/repository/user_address"
	"website-api/third-party/provider/rajaongkir"
)

type (
	IService interface {
		Find(userID string) (resData []userAddressModel.UserAddress, err error)
		Create(reqBody *userAddressModel.CreateUserAddressRequest) error
		Detail(reqPath *userAddressModel.ReqPath) (resData userAddressModel.UserAddress, err error)
		Delete(reqPath userAddressModel.ReqPath) error
		Update(reqPath *userAddressModel.UpdateUserAddress) error
		UpdateByID(req *userAddressModel.ReqUpdateUserAddress) error
	}

	service struct {
		userAddressRepo    userAddressRepo.IRepo
		masterRepo         master.IRepo
		redis              cache.Cache
		rajaOngkirProvider rajaongkir.Provider
	}
)

func NewService(userAddressRepo userAddressRepo.IRepo, masterRepo master.IRepo, redis cache.Cache, rajaOngkirProvider rajaongkir.Provider) IService {
	return &service{
		userAddressRepo:    userAddressRepo,
		masterRepo:         masterRepo,
		redis:              redis,
		rajaOngkirProvider: rajaOngkirProvider,
	}
}
