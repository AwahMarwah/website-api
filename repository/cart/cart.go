package cart

import (
	"website-api/model/cart"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(reqBody cart.Cart) error
		GetItemByUserID(userID string) (resData []cart.CartItemResponse, err error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
