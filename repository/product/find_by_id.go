package product

import productModel "website-api/model/product"

func (r *repo) FindByID(id string) (product productModel.Product, err error) {
	return product, r.db.First(&product, id).Error
}
