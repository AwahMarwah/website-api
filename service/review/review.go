package review

import (
	"errors"
	"fmt"
	"net/http"
	reviewModel "website-api/model/review"
	orderRepo "website-api/repository/order"
	"website-api/repository/review"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IService interface {
		Create(req *reviewModel.CreateReviewReq, userID string) (int, error)
		ListByProduct(productID string, page, limit, offset int) ([]reviewModel.ReviewResponse, int64, int, error)
	}

	service struct {
		reviewRepo review.IRepo
		orderRepo  orderRepo.IRepo
	}
)

func NewService(reviewRepo review.IRepo, orderRepo orderRepo.IRepo) IService {
	return &service{reviewRepo: reviewRepo, orderRepo: orderRepo}
}

func (s *service) Create(req *reviewModel.CreateReviewReq, userID string) (int, error) {
	// 1 user 1 review per produk
	if _, err := s.reviewRepo.FindByUserAndProduct(userID, req.ProductID); err == nil {
		return http.StatusConflict, fmt.Errorf("kamu sudah memberikan ulasan untuk produk ini")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, fmt.Errorf("gagal cek ulasan: %w", err)
	}

	// validasi ketat: harus punya order COMPLETED utk produk tsb
	hasOrder, err := s.orderRepo.HasCompletedOrderForProduct(userID, req.ProductID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memvalidasi pesanan: %w", err)
	}
	if !hasOrder {
		return http.StatusForbidden, fmt.Errorf("kamu hanya bisa memberikan ulasan untuk produk yang sudah selesai kamu beli")
	}

	rev := &reviewModel.Review{
		ID:        uuid.NewString(),
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	if err := s.reviewRepo.Create(rev); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menyimpan ulasan: %w", err)
	}
	return http.StatusCreated, nil
}

func (s *service) ListByProduct(productID string, page, limit, offset int) ([]reviewModel.ReviewResponse, int64, int, error) {
	res, count, err := s.reviewRepo.FindByProduct(productID, limit, offset)
	if err != nil {
		return nil, count, http.StatusInternalServerError, fmt.Errorf("gagal mengambil ulasan: %w", err)
	}
	return res, count, http.StatusOK, nil
}