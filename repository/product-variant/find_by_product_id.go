package product_variant

import product_variant "website-api/model/product-variant"

func (r *repo) FindByProductID(productID string) (productVariants []product_variant.ProductVariant, err error) {
	return productVariants, r.db.Where("product_id = ? AND is_active = ?", productID, true).Find(&productVariants).Error
}
