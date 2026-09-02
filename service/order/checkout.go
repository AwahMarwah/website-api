package order

import (
	"fmt"
	"net/http"
	"time"
	cartModel "website-api/model/cart"
	"website-api/model/order"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func (s *service) Checkout(req *order.CheckoutRequest) (order.CheckoutResponse, int, error) {
	var resData order.CheckoutResponse

	var (
		newOrder    order.Order
		totalAmount float64
		orderItems  []order.OrderItem
		variantIDs  []string
		itemDetails []midtrans.ItemDetails
	)

	err := s.txManager.Execute(func(tx *gorm.DB) error {
		txProductVariantRepo := s.productVariantRepo.WithTx(tx)

		orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

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

			variantIDs = append(variantIDs, item.VariantID)
			orderItems = append(orderItems, order.OrderItem{
				ID:               fmt.Sprintf("ORD-%d-%s", time.Now().UnixNano(), item.VariantID),
				OrderID:          orderID,
				ProductVariantID: productVariant.ID,
				Price:            productVariant.Price,
				Qty:              item.Qty,
				Subtotal:         subTotal,
			})
			itemDetails = append(itemDetails, midtrans.ItemDetails{
				ID:    item.VariantID,
				Name:  productVariant.VariantName,
				Price: int64(productVariant.Price),
				Qty:   int32(item.Qty),
			})
		}

		grandTotal := totalAmount + req.ShippingFee

		newOrder = order.Order{
			ID:            orderID,
			UserID:        req.UserID,
			AddressID:     req.AddressID,
			TotalAmount:   grandTotal,
			ShippingFee:   req.ShippingFee,
			Status:        "PENDING",
			PaymentMethod: req.PaymentMethod,
		}
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ? AND product_variant_id IN ?", req.UserID, variantIDs).Delete(&cartModel.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat order: %w", err)
	}

	// Buat Snap transaction di Midtrans setelah order berhasil dibuat
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
		// bila pembuatan payment gagal, jangan gagalkan order, tetap kembalikan error
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat transaksi pembayaran: %w", err)
	}

	// Simpan token & redirect url payment
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
