package order

import "website-api/model/order"

func toOrderResponse(o order.Order, items []order.OrderItem) order.OrderResponse {
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
		Items:         make([]order.OrderItemResponse, 0),
	}
	for _, it := range items {
		res.Items = append(res.Items, order.OrderItemResponse{
			ID:               it.ID,
			ProductVariantID: it.ProductVariantID,
			Price:            it.Price,
			Qty:              it.Qty,
			Subtotal:         it.Subtotal,
		})
	}
	return res
}
