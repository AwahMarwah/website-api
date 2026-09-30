package cart

import (
	"website-api/model/cart"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(reqBody *cart.CartItem) error
		Upsert(item *cart.CartItem) error
		SetQty(userID, variantID string, qty int) error
		GetItemByUserID(userID string) (resData []cart.CartItemResponse, err error)
		DeleteItems(userID string, variantIDs []string) error
		ClearAll(userID string) error
		WithTx(tx *gorm.DB) IRepo
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
