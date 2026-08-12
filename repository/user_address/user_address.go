package user_address

import (
	"website-api/model/user_address"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(reqBody *user_address.UserAddress) error
		Take(selectParams []string, conditions *user_address.UserAddress) (address user_address.UserAddress, err error)
		Find(userID string) (resData []user_address.UserAddress, err error)
		Delete(id string) error
		Update(id *string, values *map[string]any) (err error)
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
