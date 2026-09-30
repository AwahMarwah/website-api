package cart

import (
	"fmt"
	"net/http"
	"time"

	cartModel "website-api/model/cart"
)

// Create menambahkan item ke keranjang. Kalau varian yang sama sudah ada, qty-nya
// ditambahkan alih-alih menggagalkan karena bentrok UNIQUE constraint.
func (s *service) Create(reqBody *cartModel.CartRequestBody) (statusCode int, err error) {
	if reqBody.Qty <= 0 {
		return http.StatusBadRequest, fmt.Errorf("qty harus lebih dari 0")
	}

	item := cartModel.CartItem{
		UserID:           reqBody.UserID,
		ProductVariantID: reqBody.ProductVariantID,
		Qty:              reqBody.Qty,
		CreatedAt:        time.Now(),
	}
	if err = s.cartRepo.Upsert(&item); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}
