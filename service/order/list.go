package order

import (
	"fmt"
	"net/http"
	"website-api/model/order"
)

func (s *service) List(req *order.ListOrderReqQuery) (resData []order.OrderResponse, count int64, statusCode int, err error) {
	orders, total, err := s.orderRepo.FindByUserID(req.UserID, req.Status, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar order: %w", err)
	}

	reviewed, err := s.reviewRepo.HasReviewedMany(req.UserID, productIDsFromOrders(orders))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil status review: %w", err)
	}

	resData = make([]order.OrderResponse, 0, len(orders))
	for _, o := range orders {
		resData = append(resData, s.toOrderResponse(o, reviewed, req.UserID, true))
	}
	return resData, total, http.StatusOK, nil
}