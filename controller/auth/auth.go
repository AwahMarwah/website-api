package auth

import (
	authRepo "website-api/repository/auth"
	userRepo "website-api/repository/user"
	"website-api/service/auth"
	"website-api/worker"

	"gorm.io/gorm"
)

type controller struct {
	authService auth.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{authService: auth.NewService(userRepo.NewRepo(db), authRepo.NewRepo(db), worker.NewRedisClient())}
}
