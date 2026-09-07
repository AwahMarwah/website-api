package merchant

import (
	merchantModel "website-api/model/merchant"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		FindAll() ([]merchantModel.Merchant, error)
		FindByID(id string) (merchantModel.Merchant, error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) FindAll() ([]merchantModel.Merchant, error) {
	var merchants []merchantModel.Merchant
	err := r.db.Where("is_active = ?", true).Find(&merchants).Error
	return merchants, err
}

func (r *repo) FindByID(id string) (merchantModel.Merchant, error) {
	var m merchantModel.Merchant
	err := r.db.Where("id = ?", id).First(&m).Error
	return m, err
}