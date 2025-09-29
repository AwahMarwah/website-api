package auth

import (
	"gorm.io/gorm"
	authRepo "website-api/repository/auth"
	userRepo "website-api/repository/user"
	"website-api/service/auth"
)

type controller struct {
	authService auth.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{authService: auth.NewService(userRepo.NewRepo(db), authRepo.NewRepo(db))}
}
