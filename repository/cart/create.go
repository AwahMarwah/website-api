package cart

import "website-api/model/cart"

func (r *repo) Create(reqBody *cart.CartItem) error {
	return r.db.Create(reqBody).Error
}
