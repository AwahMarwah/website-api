package review

import (
	reviewModel "website-api/model/review"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Create(r *reviewModel.Review) error
		FindByUserAndProduct(userID, productID string) (reviewModel.Review, error)
		HasReviewed(userID, productID string) (bool, error)
		FindByProduct(productID string, limit, offset int) ([]reviewModel.ReviewResponse, int64, error)
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}

func (r *repo) Create(rev *reviewModel.Review) error {
	return r.db.Create(rev).Error
}

func (r *repo) FindByUserAndProduct(userID, productID string) (reviewModel.Review, error) {
	var rev reviewModel.Review
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&rev).Error
	return rev, err
}

func (r *repo) HasReviewed(userID, productID string) (bool, error) {
	var count int64
	err := r.db.Model(&reviewModel.Review{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repo) FindByProduct(productID string, limit, offset int) ([]reviewModel.ReviewResponse, int64, error) {
	var res []reviewModel.ReviewResponse
	var count int64

	q := r.db.Model(&reviewModel.Review{}).Where("product_id = ?", productID)
	if err := q.Count(&count).Error; err != nil {
		return nil, count, err
	}

	err := r.db.Model(&reviewModel.Review{}).
		Select("reviews.id, reviews.product_id, u.name AS user_name, reviews.rating, reviews.comment, reviews.created_at").
		Joins("JOIN users u ON u.id = reviews.user_id").
		Where("reviews.product_id = ?", productID).
		Order("reviews.created_at DESC").
		Limit(limit).Offset(offset).
		Scan(&res).Error

	return res, count, err
}