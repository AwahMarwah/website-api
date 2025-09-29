package user

import (
	"gorm.io/gorm"
	roleRepo "website-api/repository/role"
	userRepo "website-api/repository/user"
	"website-api/service/user"
)

type controller struct {
	userService user.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{userService: user.NewService(userRepo.NewRepo(db), roleRepo.NewRepo(db))}
}
