package product

import productModel "website-api/model/product"

func (r *repo) FindDetailByID(id string) (resData productModel.ProductDetailResponse, err error) {
	err = r.db.Model(&productModel.Product{}).
		Select(`products.id,
			products.name,
			products.slug,
			products.description,
			b.id AS brand_id,
			b.name AS brand_name,
			pi.image_url as thumbnail,
			products.base_price,
			MIN(pv.price) as min_price,
			MAX(pv.price) as max_price,
			MAX(CASE WHEN pv.stock > 0 THEN 1 ELSE 0 END) as is_in_stock,
			COALESCE(ROUND(AVG(r.rating),1),0) as rating,
			COUNT(DISTINCT r.id) as total_review`).
		Joins("JOIN brands b ON b.id = products.brand_id").
		Joins("JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).
		Joins("JOIN product_images pi ON pi.product_id = products.id AND pi.is_primary = ?", true).
		Joins("LEFT JOIN reviews r ON r.product_id = products.id").
		Where("products.id = ?", id).
		Group("products.id, b.id, pi.image_url, products.base_price").
		Take(&resData).Error
	return resData, err
}
