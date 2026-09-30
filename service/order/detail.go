package order

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"website-api/model/order"
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

	// Ongkir per seller ditampilkan di detail supaya pembeli tahu paket mana
	// yang dikirim merchant mana.
	shippings, err := s.orderRepo.FindMerchantShippingsByOrder(id)
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data pengiriman: %w", err)
	}
	resData.Shippings = make([]order.OrderMerchantShippingResponse, 0, len(shippings))
	for _, shp := range shippings {
		res := order.OrderMerchantShippingResponse{
			MerchantID: shp.MerchantID,
			Courier:    shp.Courier,
			Service:    shp.Service,
			Cost:       shp.Cost,
			Etd:        shp.Etd,
			WeightGram: shp.WeightGram,
		}
		if shp.Merchant != nil {
			res.MerchantName = shp.Merchant.Name
		}
		resData.Shippings = append(resData.Shippings, res)
	}

	return resData, http.StatusOK, nil
}
