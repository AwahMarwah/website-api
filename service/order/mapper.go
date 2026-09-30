package order

import "website-api/model/order"

// toOrderResponse memetakan order ke response.
//
// Nama produk memakai snapshot di order_items, bukan produk live. Kalau produknya
// diubah atau dihapus setelah transaksi, order lama tetap menampilkan nama yang
// benar-benar dibeli pembeli.
func (s *service) toOrderResponse(o order.Order, reviewed map[string]bool, userID string, checkReviewed bool) order.OrderResponse {
	res := order.OrderResponse{
		ID:              o.ID,
		AddressID:       o.AddressID,
		AddressSnapshot: o.AddressSnapshot,
		TotalAmount:     o.TotalAmount,
		ShippingFee:     o.ShippingFee,
		DiscountAmount:  o.DiscountAmount,
		Status:          o.Status,
		PaymentMethod:   o.PaymentMethod,
		PaymentURL:      o.PaymentURL,
		ExpiredAt:       o.ExpiredAt,
		Note:            o.Note,
		PaidAt:          o.PaidAt,
		ShippedAt:       o.ShippedAt,
		CompletedAt:     o.CompletedAt,
		CancelledAt:     o.CancelledAt,
		CreatedAt:       o.CreatedAt,
		Items:           make([]order.OrderItemResponse, 0, len(o.Items)),
		StatusHistories: make([]order.OrderStatusHistoryResponse, 0, len(o.StatusHistories)),
	}

	for _, it := range o.Items {
		item := order.OrderItemResponse{
			ID:               it.ID,
			ProductVariantID: it.ProductVariantID,
			MerchantID:       it.MerchantID,
			Price:            it.Price,
			Qty:              it.Qty,
			Subtotal:         it.Subtotal,
			ProductName:      it.ProductName,
			VariantName:      it.VariantName,
			ProductImageURL:  it.ProductImageURL,
			Sku:              it.Sku,
		}
		if it.Merchant != nil {
			item.MerchantName = it.Merchant.Name
		}

		if it.ProductVariant != nil {
			item.ProductID = it.ProductVariant.ProductID
			if item.ProductName == "" && it.ProductVariant.Product != nil {
				item.ProductName = it.ProductVariant.Product.Name
			}
			if checkReviewed && userID != "" {
				item.Reviewed = reviewed[it.ProductVariant.ProductID]
			}
		}

		res.Items = append(res.Items, item)
	}

	for _, h := range o.StatusHistories {
		history := order.OrderStatusHistoryResponse{
			ToStatus:  h.ToStatus,
			Note:      h.Note,
			CreatedAt: h.CreatedAt,
		}
		if h.FromStatus != nil {
			history.FromStatus = *h.FromStatus
		}
		if h.ActorRole != nil {
			history.ActorRole = *h.ActorRole
		}
		res.StatusHistories = append(res.StatusHistories, history)
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
