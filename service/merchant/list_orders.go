package merchant

import (
	"fmt"
	"net/http"

	"website-api/library/pagination"
	merchantModel "website-api/model/merchant"
	orderModel "website-api/model/order"
)

// ListOrders mengembalikan order yang memuat barang milik seller yang sedang login.
//
// Query disaring di level database lewat order_items.merchant_id yang dibekukan saat
// checkout, jadi hasilnya tidak berubah ketika seller memindahkan produk ke merchant lain.
func (s *service) ListOrders(userID string, reqQuery *merchantModel.ListSellerOrderReqQuery) ([]merchantModel.SellerOrderResponse, int64, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return nil, 0, statusCode, err
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	orders, count, err := s.orderRepo.FindByMerchantID(merchantID, reqQuery.Status, reqQuery.Limit, reqQuery.Offset)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar order: %w", err)
	}

	res := make([]merchantModel.SellerOrderResponse, 0, len(orders))
	for _, o := range orders {
		shippings, err := s.orderRepo.FindMerchantShippingsByOrder(o.ID)
		if err != nil {
			return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil data pengiriman: %w", err)
		}

		item := mapSellerOrder(o, shippings, merchantID)
		res = append(res, item)
	}
	return res, count, http.StatusOK, nil
}

// GetOrder mengembalikan detail satu order milik seller, dipangkas ke item miliknya saja.
func (s *service) GetOrder(userID, orderID string) (merchantModel.SellerOrderResponse, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return merchantModel.SellerOrderResponse{}, statusCode, err
	}

	// Ownership dicek lewat isi item, bukan lewat user_id order: satu order bisa
	// berisi barang dari beberapa seller, dan seller harus bisa melihat order itu
	// meskipun bukan pemiliknya.
	owns, err := s.orderRepo.MerchantOwnsOrder(orderID, merchantID)
	if err != nil {
		return merchantModel.SellerOrderResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal memeriksa kepemilikan order: %w", err)
	}
	if !owns {
		return merchantModel.SellerOrderResponse{}, http.StatusNotFound, fmt.Errorf("order not found")
	}

	o, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return merchantModel.SellerOrderResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil order: %w", err)
	}

	items, err := s.orderRepo.FindItemsByOrderAndMerchant(orderID, merchantID)
	if err != nil {
		return merchantModel.SellerOrderResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil item order: %w", err)
	}
	o.Items = items

	shippings, err := s.orderRepo.FindMerchantShippingsByOrder(orderID)
	if err != nil {
		return merchantModel.SellerOrderResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil data pengiriman: %w", err)
	}

	return mapSellerOrder(o, shippings, merchantID), http.StatusOK, nil
}

// mapSellerOrder memetakan order ke response seller: hanya item milik seller yang
// disertakan, dan nominal yang ditampilkan adalah bagian seller sendiri, bukan
// total seluruh order.
func mapSellerOrder(o orderModel.Order, shippings []orderModel.OrderMerchantShipping, merchantID string) merchantModel.SellerOrderResponse {
	res := merchantModel.SellerOrderResponse{
		ID:           o.ID,
		BuyerName:    o.BuyerName,
		Status:       o.Status,
		TotalAmount:  o.TotalAmount,
		CreatedAt:    o.CreatedAt,
		PaidAt:       o.PaidAt,
		ShippedAt:    o.ShippedAt,
		CompletedAt:  o.CompletedAt,
		Items:        make([]merchantModel.SellerOrderItemResponse, 0, len(o.Items)),
	}

	for _, it := range o.Items {
		if it.MerchantID != merchantID {
			continue
		}
		res.Items = append(res.Items, merchantModel.SellerOrderItemResponse{
			ID:               it.ID,
			ProductName:      it.ProductName,
			VariantName:      it.VariantName,
			ProductImageURL:  it.ProductImageURL,
			Sku:              it.Sku,
			ProductVariantID: it.ProductVariantID,
			Price:            it.Price,
			Qty:              it.Qty,
			Subtotal:         it.Subtotal,
			TotalWeightGram:  it.TotalWeightGram,
		})
		if it.ProductVariant != nil {
			res.Items[len(res.Items)-1].ProductID = it.ProductVariant.ProductID
		}
		res.ItemSubtotal += it.Subtotal
	}

	// Ongkir seller mengikuti paket yang dipilih untuk tokonya, bukan total ongkir order.
	for _, shp := range shippings {
		if shp.MerchantID == merchantID {
			res.Courier = shp.Courier
			res.Service = shp.Service
			res.Etd = shp.Etd
			res.ShippingCost = shp.Cost
			break
		}
	}
	return res
}
