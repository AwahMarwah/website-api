package order

import (
	"errors"
	"fmt"
	"net/http"
	"website-api/model/order"

	"gorm.io/gorm"
)

func (s *service) Detail(id, userID, roleName string) (resData order.OrderResponse, statusCode int, err error) {
	o, _, err := s.orderRepo.FindByIDWithItems(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusNotFound, fmt.Errorf("order not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil order: %w", err)
	}

	// Ownership check (IDOR) - hanya pemilik order atau super_admin/admin
	if o.UserID != userID && roleName != "super_admin" && roleName != "admin" {
		return resData, http.StatusForbidden, fmt.Errorf("forbidden")
	}

	checkReviewed := o.UserID == userID // hanya pemilik yang status review-nya diisi
	var reviewed map[string]bool
	if checkReviewed {
		reviewed, err = s.reviewRepo.HasReviewedMany(userID, productIDsFromOrders([]order.Order{o}))
		if err != nil {
			return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil status review: %w", err)
		}
	}
	resData = s.toOrderResponse(o, reviewed, userID, checkReviewed)
	return resData, http.StatusOK, nil
}