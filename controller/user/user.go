package user

import (
	roleRepo "website-api/repository/role"
	userRepo "website-api/repository/user"
	userAddressRepo "website-api/repository/user_address"
	"website-api/service/user"

	"gorm.io/gorm"
)

type controller struct {
	userService user.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{userService: user.NewService(userRepo.NewRepo(db), roleRepo.NewRepo(db), userAddressRepo.NewRepo(db))}
}
