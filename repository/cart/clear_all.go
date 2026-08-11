package cart

import cartModel "website-api/model/cart"

func (r *repo) ClearAll(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&cartModel.CartItem{}).Error
}
