package master

import (
	masterRepo "website-api/repository/master"
	"website-api/service/master"

	"gorm.io/gorm"
)

type controller struct {
	masterService master.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{masterService: master.NewService(masterRepo.NewRepo(db))}
}
