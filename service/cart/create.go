package cart

import (
	"net/http"
	"time"
	cartModel "website-api/model/cart"
)

func (s *service) Create(reqBody *cartModel.CartRequestBody) (statusCode int, err error) {
	item := cartModel.CartItem{
		UserID:           reqBody.UserID,
		ProductVariantID: reqBody.ProductVariantID,
		Qty:              reqBody.Qty,
		CreatedAt:        time.Now(),
	}
	if err = s.cartRepo.Create(&item); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}
