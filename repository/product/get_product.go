package product

import (
	"website-api/library/helper/filter"
	productModel "website-api/model/product"
)

func (r *repo) GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error) {
	resData = make([]productModel.ListProductResponse, 0)

	err = r.db.Debug().Model(&productModel.Product{}).
		Select(`products.id, 
			products.name, 
			products.slug, 
			b.id AS brand_id,
			b.name AS brand_name,
			pi.image_url as thumbnail,
			MIN(pv.price) as min_price,
			MAX(pv.price) as max_price,
			MAX(CASE WHEN pv.stock > 0 THEN 1 ELSE 0 END) as is_in_stock,
			COALESCE(ROUND(AVG(r.rating),1),0) as rating,
			COUNT(DISTINCT r.id) as total_review`).
		Joins("JOIN brands b ON b.id = products.brand_id").
		Joins("JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).
		Joins("JOIN product_images pi ON pi.product_id = products.id AND pi.is_primary = ?", true).
		Joins("LEFT JOIN reviews r ON r.product_id = products.id").
		Joins("LEFT JOIN product_categories pc ON pc.product_id = products.id").
		Joins("LEFT JOIN categories c ON c.id = pc.category_id").
		Scopes(
			filter.FilterBrand(reqQuery.Brand),
			filter.FilterCategory(reqQuery.Category),
			filter.FilterMaxPrice(reqQuery.MaxPrice),
			filter.FilterMinPrice(reqQuery.MinPrice),
		).
		Group("products.id, b.id, pi.image_url").
		Count(&count).Limit(reqQuery.Limit).Offset(reqQuery.Offset).Find(&resData).Error
	return resData, count, nil
}
