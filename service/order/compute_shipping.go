package order

import (
	"fmt"
	"math"

	"website-api/model/order"
	"website-api/third-party/provider/rajaongkir"
)

// computeShippingBreakdown menghitung ongkir per merchant berdasarkan shippings yang dipilih user.
// Berat dipisah per merchant karena ongkir dihitung dari kota asal seller, bukan satu asal untuk
// seluruh order.
func (s *service) computeShippingBreakdown(items []order.CheckoutItem, destinationID int64, shippings []order.ShippingsReq) ([]order.OrderMerchantShipping, float64, error) {
	weightByMerchant := map[string]int{}
	merchantByID := map[string]string{} // merchantID -> name

	for _, item := range items {
		variant, err := s.productVariantRepo.FindByID(item.VariantID)
		if err != nil {
			return nil, 0, fmt.Errorf("variant %s not found", item.VariantID)
		}
		snapshot, err := s.productRepo.FindCheckoutSnapshot(variant.ProductID)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mengambil produk: %w", err)
		}
		if snapshot.MerchantID == "" {
			return nil, 0, fmt.Errorf("produk %s belum memiliki merchant", snapshot.ProductName)
		}
		merchantByID[snapshot.MerchantID] = snapshot.ProductName
		weightGramPerUnit := int(math.Round(float64(variant.Weight) * 1000))
		weightByMerchant[snapshot.MerchantID] += weightGramPerUnit * item.Qty
	}

	result := make([]order.OrderMerchantShipping, 0, len(weightByMerchant))
	var total float64 = 0

	for merchantID, weightGram := range weightByMerchant {
		var selected *order.ShippingsReq
		for i := range shippings {
			if shippings[i].MerchantID == merchantID {
				selected = &shippings[i]
				break
			}
		}
		if selected == nil {
			return nil, 0, fmt.Errorf("shipping untuk merchant %s belum dipilih", merchantByID[merchantID])
		}

		merchant, err := s.merchantRepo.FindByID(merchantID)
		if err != nil {
			return nil, 0, fmt.Errorf("merchant %s not found", merchantID)
		}

		options, err := s.rajaOngkir.CalculateCost(rajaongkir.CalculateCostRequest{
			Origin:      merchant.DestinationID,
			Destination: destinationID,
			Weight:      weightGram,
			Courier:     selected.Courier,
		})
		if err != nil {
			return nil, 0, fmt.Errorf("gagal menghitung ongkir merchant %s: %w", merchant.Name, err)
		}

		cost := int64(0)
		etd := ""
		found := false
		for _, opt := range options {
			if opt.Service == selected.Service {
				cost = opt.Cost
				etd = opt.Etd
				found = true
				break
			}
		}
		if !found {
			return nil, 0, fmt.Errorf("kurir/service %s - %s tidak tersedia untuk merchant %s", selected.Courier, selected.Service, merchant.Name)
		}

		result = append(result, order.OrderMerchantShipping{
			MerchantID: merchantID,
			Courier:    selected.Courier,
			Service:    selected.Service,
			Cost:       cost,
			Etd:        etd,
			WeightGram: weightGram,
		})
		total += float64(cost)
	}

	return result, total, nil
}
