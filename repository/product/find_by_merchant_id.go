package product

import (
	"fmt"
	"strings"

	productModel "website-api/model/product"

	"gorm.io/gorm"
)

// FindByMerchantID mengikuti pola GetProduct: daftar produk dirakit lewat join agar
// thumbnail, rentang harga, dan agregat rating bisa diambil dalam satu query.
func (r *repo) FindByMerchantID(merchantID, status, search, voucherID, sort string, limit, offset int) (resData []productModel.ListProductResponse, count int64, err error) {
	resData = make([]productModel.ListProductResponse, 0)

	scopes := []func(db *gorm.DB) *gorm.DB{
		filterMerchant(merchantID),
		filterProductStatus(status),
		filterMerchantProductSearch(search),
		filterVoucherProducts(voucherID),
	}

	countQuery := r.db.Model(&productModel.Product{}).
		Joins("JOIN brands b ON b.id = products.brand_id").
		Joins("JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).
		Scopes(scopes...).
		Distinct("products.id")
	if err = countQuery.Count(&count).Error; err != nil {
		return nil, count, fmt.Errorf("gagal menghitung produk: %w", err)
	}

	q := r.db.Model(&productModel.Product{}).
		Select(`products.id,
			products.name,
			products.slug,
			products.status,
			b.id AS brand_id,
			b.name AS brand_name,
			m.id AS merchant_id,
			m.name AS merchant_name,
			COALESCE(pi.image_url, '') as thumbnail,
			COALESCE(MIN(pv.price), 0) as min_price,
			COALESCE(MAX(pv.price), 0) as max_price,
			MAX(CASE WHEN pv.stock > 0 THEN 1 ELSE 0 END) as is_in_stock,
			COALESCE(ROUND(AVG(r.rating),1),0) as rating,
			COUNT(DISTINCT r.id) as total_review`).
		Joins("JOIN brands b ON b.id = products.brand_id").
		Joins("JOIN merchants m ON m.id = products.merchant_id").
		Joins("JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).
		Joins("LEFT JOIN product_images pi ON pi.product_id = products.id AND pi.is_primary = ?", true).
		Joins("LEFT JOIN reviews r ON r.product_id = products.id").
		Scopes(scopes...).
		Group("products.id, b.id, m.id, pi.image_url")

	switch strings.ToLower(sort) {
	case "price_asc":
		q = q.Order("min_price ASC")
	case "price_desc":
		q = q.Order("max_price DESC")
	case "rating":
		q = q.Order("rating DESC")
	default:
		q = q.Order("products.created_at DESC")
	}

	if err = q.Limit(limit).Offset(offset).Find(&resData).Error; err != nil {
		return nil, count, fmt.Errorf("gagal mengambil daftar produk: %w", err)
	}
	return resData, count, nil
}

func filterMerchant(merchantID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("products.merchant_id = ?", merchantID)
	}
}

func filterProductStatus(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status == "" {
			return db
		}
		return db.Where("products.status = ?", status)
	}
}

func filterMerchantProductSearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		return db.Where("products.name ILIKE ? OR products.sku ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
}

// filterVoucherProducts membatasi ke produk yang terdaftar pada voucher. Voucher
// tanpa daftar produk berarti berlaku untuk seluruh produk seller, jadi tidak difilter.
func filterVoucherProducts(voucherID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if voucherID == "" {
			return db
		}
		return db.Where("products.id IN (?)", db.Session(&gorm.Session{NewDB: true}).
			Table("product_vouchers").
			Select("product_id").
			Where("voucher_id = ?", voucherID))
	}
}
