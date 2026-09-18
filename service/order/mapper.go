package order

import (
	"website-api/model/order"
)

func (s *service) toOrderResponse(o order.Order, items []order.OrderItem, userID string, checkReviewed bool) (order.OrderResponse, error) {
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
		item := order.OrderItemResponse{
			ID:               it.ID,
			ProductVariantID: it.ProductVariantID,
			Price:            it.Price,
			Qty:              it.Qty,
			Subtotal:         it.Subtotal,
		}

		// resolve variant → product id & name
		variant, err := s.productVariantRepo.FindByID(it.ProductVariantID)
		if err == nil {
			item.ProductID = variant.ProductID
			if p, err := s.productRepo.FindByID(variant.ProductID); err == nil {
				item.ProductName = p.Name
			}
			if checkReviewed && userID != "" {
				reviewed, rErr := s.reviewRepo.HasReviewed(userID, variant.ProductID)
				if rErr == nil {
					item.Reviewed = reviewed
				}
			}
		}

		res.Items = append(res.Items, item)
	}
	return res, nil
}