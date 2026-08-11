package order

import (
	"fmt"
	"time"
	cartModel "website-api/model/cart"
	"website-api/model/order"

	"gorm.io/gorm"
)

func (s *service) Checkout(req *order.CheckoutRequest) error {
	return s.txManager.Execute(func(tx *gorm.DB) error {
		//txProductRepo := s.productRepo.WithTx(tx)
		txProductVariantRepo := s.productVariantRepo.WithTx(tx)

		var (
			totalAmount float64
			orderItems  []order.OrderItem
			variantIDs  []string
		)

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

			// validasi stock variant
			if productVariant.Stock < item.Qty {
				return fmt.Errorf("stock not enough for variant %s (%s)", productVariant.VariantName, productVariant.Sku)
			}

			// kurangi stock varian
			productVariant.Stock -= item.Qty
			if err := txProductVariantRepo.Update(productVariant); err != nil {
				return err
			}

			// hitung subtotal & total
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
		}
		grandTotal := totalAmount + req.ShippingFee

		newOrder := order.Order{
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
}
