package cart

import cartModel "website-api/model/cart"

func (r *repo) DeleteItems(userID string, variantIDs []string) error {
	return r.db.Where("user_id = ? AND product_variant_id IN ?", userID, variantIDs).Delete(&cartModel.CartItem{}).Error
}
