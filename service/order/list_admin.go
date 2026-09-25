package order

import (
	"fmt"
	"net/http"
	"website-api/model/order"
)

func (s *service) ListAdmin(req *order.ListOrderReqQuery) (resData []order.OrderResponse, count int64, statusCode int, err error) {
	orders, total, err := s.orderRepo.FindAll(req.Status, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar order: %w", err)
	}

	resData = make([]order.OrderResponse, 0, len(orders))
	for _, o := range orders {
		resData = append(resData, s.toOrderResponse(o, nil, "", false))
	}
	return resData, total, http.StatusOK, nil
}