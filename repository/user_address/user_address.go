package user_address

import (
	"website-api/model/user_address"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Take(selectParams []string, conditions *user_address.UserAddress) (address user_address.UserAddress, err error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{
		db: db,
	}
}
