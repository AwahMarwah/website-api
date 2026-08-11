package cart

import (
	"fmt"
	"net/http"
	cartModel "website-api/model/cart"
	"website-api/utils"
)

func (s *service) GetCartByUserID(userID string) (resData cartModel.GetCartByUserIDResponse, statusCode int, err error) {
	items, err := s.cartRepo.GetItemByUserID(userID)
	fmt.Println(items, "item service")
	if err != nil {
		statusCode = http.StatusInternalServerError
	}
	var grandTotal float64
	for i := range items {
		grandTotal += items[i].SubTotal
		items[i].FormattedPrice = utils.FormatRupiah(items[i].Price)
		items[i].FormattedSubTotal = utils.FormatRupiah(items[i].SubTotal)
	}

	resData.UserID = userID
	resData.Items = items
	resData.Summary.TotalItems = len(items)
	resData.Summary.GrandTotal = grandTotal
	resData.Summary.FormattedGrandTotal = utils.FormatRupiah(grandTotal)
	return resData, statusCode, nil
}
