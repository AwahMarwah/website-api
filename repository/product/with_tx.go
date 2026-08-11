package product

import "gorm.io/gorm"

func (r *repo) WithTx(tx *gorm.DB) IRepo {
	if tx == nil {
		return r
	}
	return &repo{
		db: tx,
	}
}
