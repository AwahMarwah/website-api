package order

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"website-api/common"
	cartModel "website-api/model/cart"
	"website-api/model/order"
	userModel "website-api/model/user"
	userAddressModel "website-api/model/user_address"

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
		discount    float64
		orderItems  []order.OrderItem
		variantIDs  []string
		itemDetails []midtrans.ItemDetails
		addressSnap order.AddressSnapshot
		buyerName   string
		buyerEmail  string
		buyerPhone  string
	)

	// Ambil alamat tujuan + ownership
	address, err := s.userAddressRepo.Take(
		[]string{"id", "user_id", "destination_id", "recipient_name", "phone_number",
			"full_address", "city", "postal_code", "province_id", "city_id", "district_id", "subdistrict_id"},
		&userAddressModel.UserAddress{ID: req.AddressID},
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, http.StatusBadRequest, fmt.Errorf("address not found")
		}
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal mengambil alamat: %w", err)
	}
	if address.UserID != req.UserID {
		return resData, http.StatusForbidden, fmt.Errorf("forbidden")
	}

	// Bekukan alamat dan identitas penerima. Setelah ini, perubahan atau penghapusan
	// alamat oleh pengguna tidak boleh mengubah tampilan order lama.
	addressSnap = order.AddressSnapshot{
		RecipientName: address.RecipientName,
		PhoneNumber:   address.PhoneNumber,
		FullAddress:   address.FullAddress,
		City:          address.City,
		PostalCode:    address.PostalCode,
		ProvinceID:    address.ProvinceID,
		CityID:        address.CityID,
		DistrictID:    address.DistrictID,
		SubdistrictID: address.SubdistrictID,
		DestinationID: address.DestinationID,
	}

	if buyer, err := s.userRepo.Take([]string{"id", "name", "email", "phone_number"}, &userModel.User{Id: req.UserID}); err == nil {
		buyerName = buyer.Name
		buyerEmail = buyer.Email
		buyerPhone = buyer.PhoneNumber
	}
	if buyerPhone == "" {
		buyerPhone = address.PhoneNumber
	}
	if buyerName == "" {
		buyerName = address.RecipientName
	}

	if len(req.Items) == 0 {
		return resData, http.StatusBadRequest, fmt.Errorf("keranjang kosong")
	}

	// Hitung ongkir per merchant bila request menyertakan shippings
	if len(req.Shippings) > 0 {
		if address.DestinationID == 0 {
			return resData, http.StatusBadRequest, fmt.Errorf("alamat belum memiliki destination_id")
		}
		shipments, totalShip, err = s.computeShippingBreakdown(req.Items, address.DestinationID, req.Shippings)
		if err != nil {
			return resData, http.StatusBadRequest, err
		}
	} else {
		// Fallback (tanpa shippings): pakai shipping_fee dari client.
		totalShip = req.ShippingFee
	}

	orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())

	var appliedVoucher *voucherRedemption

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		txProductVariantRepo := s.productVariantRepo.WithTx(tx)
		txOrderRepo := s.orderRepo.WithTx(tx)

		// Kunci baris variant satu per satu, dalam urutan ID yang stabil supaya
		// dua checkout bersamaan tidak saling menunggu (deadlock).
		sortedItems := sortedByVariantID(req.Items)

		for _, item := range sortedItems {
			// SELECT ... FOR UPDATE: baris terkunci sampai transaksi selesai, jadi
			// pengecekan stok dan pengurangan stok selalu melihat nilai yang sama.
			productVariant, err := txProductVariantRepo.FindByIDForUpdate(item.VariantID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
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

			// Kunci pemilik barang dan data tampilan saat transaksi terjadi.
			// Inilah yang membuat seller bisa melihat order-nya tanpa join ke
			// tabel produk, dan menjaga histori tetap utuh walau produk berubah.
			snapshot, err := s.productRepo.FindCheckoutSnapshot(productVariant.ProductID)
			if err != nil {
				return fmt.Errorf("gagal mengambil data produk: %w", err)
			}

			subTotal := float64(item.Qty) * productVariant.Price
			totalAmount += subTotal
			weightGramPerUnit := int(math.Round(float64(productVariant.Weight) * 1000))

			variantIDs = append(variantIDs, item.VariantID)
			orderItems = append(orderItems, order.OrderItem{
				ID:               fmt.Sprintf("ORD-%d-%s", time.Now().UnixNano(), item.VariantID),
				OrderID:          orderID,
				ProductVariantID: productVariant.ID,
				MerchantID:       snapshot.MerchantID,
				Price:            productVariant.Price,
				Qty:              item.Qty,
				Subtotal:         subTotal,
				ProductName:      snapshot.ProductName,
				VariantName:      productVariant.VariantName,
				ProductImageURL:  snapshot.ImageURL,
				Sku:              productVariant.Sku,
				TotalWeightGram:  weightGramPerUnit * item.Qty,
			})
			itemDetails = append(itemDetails, midtrans.ItemDetails{
				ID:    item.VariantID,
				Name:  productVariant.VariantName,
				Price: int64(productVariant.Price),
				Qty:   int32(item.Qty),
			})
		}

		// Validasi voucher dilakukan di sini, bukan sebelum transaksi dibuka:
		// subtotal baru diketahui setelah harga dibaca dari variant yang terkunci,
		// sedangkan voucher punya syarat minimum belanja. Memvalidasinya lebih awal
		// berarti sproket dihitung terhadap ongkir saja.
		if req.VoucherCode != "" {
			redemption, voucherErr := s.validateVoucher(req.VoucherCode, req.UserID, totalAmount+totalShip, req.Items)
			if voucherErr != nil {
				return voucherErr
			}
			appliedVoucher = redemption
			// Diskon disimpan sebagai nominal jadi, bukan dihitung ulang saat
			// ditampilkan, agar voucher yang nanti mati tidak mengubah order lama.
			discount = appliedVoucher.Discount
		}

		grandTotal := totalAmount + totalShip - discount
		if grandTotal < 0 {
			grandTotal = 0
		}

		newOrder = order.Order{
			ID:              orderID,
			UserID:          req.UserID,
			AddressID:       req.AddressID,
			AddressSnapshot: addressSnap,
			BuyerName:       buyerName,
			BuyerEmail:      buyerEmail,
			BuyerPhone:      buyerPhone,
			Note:            optionalString(req.Note),
			TotalAmount:     grandTotal,
			ShippingFee:     totalShip,
			DiscountAmount:  discount,
			Status:          common.OrderStatusPending,
			PaymentMethod:   req.PaymentMethod,
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

		// Voucher diklaim di akhir transaksi. Kuotanya bertambah di dalam transaksi
		// yang sama, jadi kalau ada langkah berikutnya yang gagal, penambahan kuota
		// ikut ter-rollback.
		if appliedVoucher != nil {
			if err := s.redeemVoucher(tx, appliedVoucher, orderID, req.UserID, discount); err != nil {
				return err
			}
		}

		if err := tx.Where("user_id = ? AND product_variant_id IN ?", req.UserID, variantIDs).Delete(&cartModel.CartItem{}).Error; err != nil {
			return err
		}

		fromStatus := ""
		return txOrderRepo.CreateStatusHistory(order.OrderStatusHistory{
			OrderID:    orderID,
			FromStatus: &fromStatus,
			ToStatus:   common.OrderStatusPending,
			ActorID:    &req.UserID,
			Note:       optionalString("order dibuat"),
			CreatedAt:  time.Now(),
		})
	})
	if err != nil {
		return resData, http.StatusInternalServerError, fmt.Errorf("gagal membuat order: %w", err)
	}

	// Snapshot transaksi Midtrans
	expiresAt := time.Now().Add(24 * time.Hour)
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newOrder.ID,
			GrossAmt: int64(math.Round(newOrder.TotalAmount)),
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
		// Transaksi pembayaran gagal dibuat. Order tetap tercatat supaya bisa
		// dibatalkan dan stokenya dikembalikan lewat cron kedaluwarsa.
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
		TotalAmount:  newOrder.TotalAmount,
	}
	if snapResp.RedirectURL != "" {
		url := snapResp.RedirectURL
		resData.PaymentURL = &url
	}

	return resData, http.StatusOK, nil
}

// sortedByVariantID mengurutkan item berdasarkan variant id agar urutan penguncian
// baris konsisten di semua transaksi. Tanpa ini, dua checkout yang berisi varian
// yang sama dalam urutan berbeda bisa mengunci baris dalam urutan berbeda dan deadlock.
func sortedByVariantID(items []order.CheckoutItem) []order.CheckoutItem {
	sorted := make([]order.CheckoutItem, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].VariantID < sorted[j].VariantID
	})
	return sorted
}
