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

	resData = make([]order.OrderResponse, 0)
	for _, o := range orders {
		items, err := s.orderRepo.FindItemsByOrderID(o.ID)
		if err != nil {
			return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil item order: %w", err)
		}
		resData = append(resData, toOrderResponse(o, items))
	}
	return resData, total, http.StatusOK, nil
}
