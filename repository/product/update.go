package product

import productModel "website-api/model/product"

func (r *repo) Update(product productModel.Product) error {
	return r.db.Save(product).Error
}
