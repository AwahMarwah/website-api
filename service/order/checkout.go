package order

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
	cartModel "website-api/model/cart"
	"website-api/model/order"
	userAddressModel "website-api/model/user_address"
	"website-api/third-party/provider/rajaongkir"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func (s *service) Checkout(req *order.CheckoutRequest) (order.CheckoutResponse, int, error) {
	var resData order.CheckoutResponse

	var (
		newOrder    order.Order
		shipments   []order.OrderMerchantShipping
		totalAmount float64
		totalShip   float64
		orderItems  []order.OrderItem
		variantIDs  []string
		itemDetails []midtrans.ItemDetails
	)

	// Ambil alamat tujuan + ownership
	address, err := s.userAddressRepo.Take([]string{"id", "user_id", "destination_id"}, &userAddressModel.UserAddress{ID: req.AddressID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusBadRequest, fmt.Errorf("address not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil alamat: %w", err)
	}
	if address.UserID != req.UserID {
		return resData, http.StatusForbidden, fmt.Errorf("forbidden")
	}

	// Hitung ongkir per merchant bila request menyertakan shippings
	if len(req.Shippings) > 0 {
		if address.DestinationID == 0 {
			return resData, http.StatusBadRequest, fmt.Errorf("alamat belum memiliki destination_id")
		}
		shipments, totalShip, err = s.computeShippingBreakdown(req.Items, address.DestinationID, req.Shippings)
		if err != nil {
			return resData, http.StatusInternalServerError, err
		}
	} else {
		// fallback (tanpa shippings): pakai shipping_fee dari client
		totalShip = req.ShippingFee
	}

	orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		txProductVariantRepo := s.productVariantRepo.WithTx(tx)
		txOrderRepo := s.orderRepo.WithTx(tx)

		for _, item := range req.Items {
			productVariant, err := txProductVariantRepo.FindByID(item.VariantID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("variant with id %s not found", item.VariantID)
				}
				return err
			}

			if !productVariant.IsActive {
				return fmt.Errorf("variant %s is inactive", productVariant.VariantName)
			}
			if productVariant.Stock < item.Qty {
				return fmt.Errorf("stock not enough for variant %s (%s)", productVariant.VariantName, productVariant.Sku)
			}

			productVariant.Stock -= item.Qty
			if err := txProductVariantRepo.Update(productVariant); err != nil {
				return err
			}

			subTotal := float64(item.Qty) * productVariant.Price
			totalAmount += subTotal
			weightGramPerUnit := int(math.Round(float64(productVariant.Weight) * 1000))

			variantIDs = append(variantIDs, item.VariantID)
			orderItems = append(orderItems, order.OrderItem{
				ID:               fmt.Sprintf("ORD-%d-%s", time.Now().UnixNano(), item.VariantID),
				OrderID:          orderID,
				ProductVariantID: productVariant.ID,
				Price:            productVariant.Price,
				Qty:              item.Qty,
				Subtotal:         subTotal,
				TotalWeightGram:  weightGramPerUnit * item.Qty,
			})
			itemDetails = append(itemDetails, midtrans.ItemDetails{
				ID:    item.VariantID,
				Name:  productVariant.VariantName,
				Price: int64(productVariant.Price),
				Qty:   int32(item.Qty),
			})
		}

		grandTotal := totalAmount + totalShip

		newOrder = order.Order{
			ID:            orderID,
			UserID:        req.UserID,
			AddressID:     req.AddressID,
			TotalAmount:   grandTotal,
			ShippingFee:   totalShip,
			Status:        "PENDING",
			PaymentMethod: req.PaymentMethod,
		}
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// Simpan breakdown ongkir per merchant
		for _, shp := range shipments {
			shp.ID = uuid.NewString()
			shp.OrderID = orderID
			if err := txOrderRepo.CreateMerchantShipping(&shp); err != nil {
				return err
			}
		}

		if err := tx.Where("user_id = ? AND product_variant_id IN ?", req.UserID, variantIDs).Delete(&cartModel.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat order: %w", err)
	}

	// Snap transaction
	expiresAt := time.Now().Add(24 * time.Hour)
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newOrder.ID,
			GrossAmt: int64(newOrder.TotalAmount),
		},
		Items: &itemDetails,
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Unit:      "hours",
			Duration:  24,
		},
	}

	snapResp, err := s.midtransProvider.CreateTransaction(snapReq)
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat transaksi pembayaran: %w", err)
	}

	updateMap := map[string]interface{}{
		"payment_token": snapResp.Token,
		"expired_at":    expiresAt,
	}
	if snapResp.RedirectURL != "" {
		updateMap["payment_url"] = snapResp.RedirectURL
	}
	if err := s.orderRepo.UpdatePaymentInfo(newOrder.ID, updateMap); err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal menyimpan info pembayaran: %w", err)
	}

	resData = order.CheckoutResponse{
		OrderID:      newOrder.ID,
		PaymentToken: snapResp.Token,
		ExpiresAt:    &expiresAt,
	}
	if snapResp.RedirectURL != "" {
		url := snapResp.RedirectURL
		resData.PaymentURL = &url
	}

	return resData, http.StatusOK, nil
}

// computeShippingBreakdown menghitung ongkir per merchant berdasarkan shippings yang dipilih user.
func (s *service) computeShippingBreakdown(items []order.CheckoutItem, destinationID int64, shippings []order.ShippingsReq) ([]order.OrderMerchantShipping, float64, error) {
	weightByMerchant := map[string]int{}
	merchantByID := map[string]string{} // merchantID -> name

	for _, item := range items {
		variant, err := s.productVariantRepo.FindByID(item.VariantID)
		if err != nil {
			return nil, 0, fmt.Errorf("variant %s not found", item.VariantID)
		}
		productInfo, err := s.productRepo.FindByID(variant.ProductID)
		if err != nil {
			return nil, 0, fmt.Errorf("gagal mengambil produk: %w", err)
		}
		if productInfo.MerchantId == "" {
			return nil, 0, fmt.Errorf("produk %s belum memiliki merchant", productInfo.Name)
		}
		if _, ok := merchantByID[productInfo.MerchantId]; !ok {
			m, err := s.merchantRepo.FindByID(productInfo.MerchantId)
			if err != nil {
				return nil, 0, fmt.Errorf("merchant %s not found", productInfo.MerchantId)
			}
			merchantByID[m.ID] = m.Name
		}
		weightGramPerUnit := int(math.Round(float64(variant.Weight) * 1000))
		weightByMerchant[productInfo.MerchantId] += weightGramPerUnit * item.Qty
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