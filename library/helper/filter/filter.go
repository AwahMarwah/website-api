package filter

import "gorm.io/gorm"

func FilterBrand(brand string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if brand != "" {
			return db.Where("b.slug = ?", brand)
		}
		return db
	}
}

func FilterCategory(category string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if category != "" {
			return db.Where("c.slug = ?", category)
		}
		return db
	}
}

func FilterMinPrice(min float64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if min > 0 {
			return db.Where("pv.price >= ?", min)
		}
		return db
	}
}

func FilterMaxPrice(max float64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if max > 0 {
			return db.Where("pv.price <= ?", max)
		}
		return db
	}
}
