package cart

import "website-api/model/cart"

func (r *repo) GetItemByUserID(userID string) (resData []cart.CartItemResponse, err error) {
	selectQuery := `
		ci.id AS cart_item_id,
		p.id AS product_id,
		p.name AS product_name,
		pv.id AS product_variant_id,
		pv.variant_name AS variant_name,
		pv.price AS price,
		ci.qty AS qty,
		(ci.qty * pv.price) AS sub_total,
		pv.stock AS current_stock`

	err = r.db.
		Table("cart_items AS ci").
		Select(selectQuery).
		Joins("INNER JOIN product_variants AS pv ON pv.id = ci.product_variant_id").
		Joins("INNER JOIN products AS p ON p.id = pv.product_id").
		Where("ci.user_id = ?", userID).
		Scan(&resData).Error

	if err != nil {
		return nil, err
	}

	return resData, nil
}
