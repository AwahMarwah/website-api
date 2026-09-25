package product

import (
	"fmt"
	"strings"
	"website-api/library/helper/filter"
	productModel "website-api/model/product"

	"gorm.io/gorm"
)

func (r *repo) GetProduct(reqQuery *productModel.GetListProductReqQuerry) (resData []productModel.ListProductResponse, count int64, err error) {
	resData = make([]productModel.ListProductResponse, 0)

	scopes := []func(db *gorm.DB) *gorm.DB{
		filter.FilterBrand(reqQuery.Brand),
		filter.FilterCategory(reqQuery.Category),
		filter.FilterMaxPrice(reqQuery.MaxPrice),
		filter.FilterMinPrice(reqQuery.MinPrice),
		filter.FilterProductSearch(reqQuery.Search),
	}

	// Hitung total produk yang cocok filter (distinct product id agar tidak
	// terinflasi oleh join multiple-row seperti variant/kategori).
	countQuery := r.db.Model(&productModel.Product{}).
		Joins("JOIN brands b ON b.id = products.brand_id").
		Joins("JOIN merchants m ON m.id = products.merchant_id").
		Joins("JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).
		Joins("LEFT JOIN product_categories pc ON pc.product_id = products.id").
		Joins("LEFT JOIN categories c ON c.id = pc.category_id").
		Scopes(scopes...).
		Distinct("products.id")
	if err = countQuery.Count(&count).Error; err != nil {
		return nil, count, fmt.Errorf("gagal menghitung produk: %w", err)
	}

	q := r.db.Model(&productModel.Product{}).
		Select(`products.id,
			products.name,
			products.slug,
			b.id AS brand_id,
			b.name AS brand_name,
			m.id AS merchant_id,
			m.name AS merchant_name,
			pi.image_url as thumbnail,
			MIN(pv.price) as min_price,
			MAX(pv.price) as max_price,
			MAX(CASE WHEN pv.stock > 0 THEN 1 ELSE 0 END) as is_in_stock,
			COALESCE(ROUND(AVG(r.rating),1),0) as rating,
			COUNT(DISTINCT r.id) as total_review`).
		Joins("JOIN brands b ON b.id = products.brand_id").
		Joins("JOIN merchants m ON m.id = products.merchant_id").
		Joins("JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).
		Joins("LEFT JOIN product_images pi ON pi.product_id = products.id AND pi.is_primary = ?", true).
		Joins("LEFT JOIN reviews r ON r.product_id = products.id").
		Joins("LEFT JOIN product_categories pc ON pc.product_id = products.id").
		Joins("LEFT JOIN categories c ON c.id = pc.category_id").
		Scopes(scopes...).
		Group("products.id, b.id, m.id, pi.image_url")

	// Sorting
	switch strings.ToLower(reqQuery.Sort) {
	case "price_asc":
		q = q.Order("min_price ASC")
	case "price_desc":
		q = q.Order("max_price DESC")
	case "rating":
		q = q.Order("rating DESC")
	default:
		q = q.Order("products.created_at DESC")
	}

	if err = q.Limit(reqQuery.Limit).Offset(reqQuery.Offset).Find(&resData).Error; err != nil {
		return nil, count, fmt.Errorf("gagal mengambil daftar produk: %w", err)
	}
	return resData, count, nil
}