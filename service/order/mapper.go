package order

import "website-api/model/order"

func (s *service) toOrderResponse(o order.Order, reviewed map[string]bool, userID string, checkReviewed bool) order.OrderResponse {
	res := order.OrderResponse{
		ID:            o.ID,
		AddressID:     o.AddressID,
		TotalAmount:   o.TotalAmount,
		ShippingFee:   o.ShippingFee,
		Status:        o.Status,
		PaymentMethod: o.PaymentMethod,
		PaymentURL:    o.PaymentURL,
		ExpiredAt:     o.ExpiredAt,
		CreatedAt:     o.CreatedAt,
		Items:         make([]order.OrderItemResponse, 0, len(o.Items)),
	}

	for _, it := range o.Items {
		item := order.OrderItemResponse{
			ID:               it.ID,
			ProductVariantID: it.ProductVariantID,
			Price:            it.Price,
			Qty:              it.Qty,
			Subtotal:         it.Subtotal,
		}

		if it.ProductVariant != nil {
			item.ProductID = it.ProductVariant.ProductID
			if it.ProductVariant.Product != nil {
				item.ProductName = it.ProductVariant.Product.Name
			}
			if checkReviewed && userID != "" {
				item.Reviewed = reviewed[it.ProductVariant.ProductID]
			}
		}

		res.Items = append(res.Items, item)
	}
	return res
}

func productIDsFromOrders(orders []order.Order) []string {
	seen := make(map[string]struct{})
	var ids []string
	for _, o := range orders {
		for _, it := range o.Items {
			if it.ProductVariant == nil {
				continue
			}
			productID := it.ProductVariant.ProductID
			if _, ok := seen[productID]; !ok {
				seen[productID] = struct{}{}
				ids = append(ids, productID)
			}
		}
	}
	return ids
}