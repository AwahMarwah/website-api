package cart

import (
	cartModel "website-api/model/cart"
	"website-api/repository/cart"
)

type (
	IService interface {
		Create(reqBody *cartModel.CartRequestBody) (statusCode int, err error)
		GetCartByUserID(userID string) (resData cartModel.GetCartByUserIDResponse, statusCode int, err error)
		DeleteItems(req *cartModel.DeleteCartRequest) error
	}

	service struct {
		cartRepo cart.IRepo
	}
)

func NewService(cartRepo cart.IRepo) IService {
	return &service{cartRepo: cartRepo}
}
