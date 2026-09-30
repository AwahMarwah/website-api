package order

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"website-api/model/order"
	voucherModel "website-api/model/voucher"
)

// voucherRedemption adalah voucher yang lolos validasi dan siap diklaim saat checkout.
type voucherRedemption struct {
	Voucher  voucherModel.Voucher
	Discount float64
}

// validateVoucher mengecek kode voucher setelah subtotal dan ongkir sudah diketahui.
//
// Syarat non-database delegate ke Voucher.Validate di model, supaya aturan yang sama
// dipakai oleh pratinjau di admin. Syarat yang butuh database dicek di sini.
func (s *service) validateVoucher(code, userID string, amount float64, items []order.CheckoutItem) (*voucherRedemption, error) {
	v, err := s.voucherRepo.FindByCode(strings.ToUpper(strings.TrimSpace(code)))
	if err != nil {
		return nil, voucherModel.NewValidationError("kode voucher tidak ditemukan")
	}

	if err := v.Validate(time.Now(), amount); err != nil {
		return nil, err
	}

	if v.PerUserLimit > 0 {
		used, err := s.voucherRepo.CountUserRedemptions(v.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("gagal memeriksa riwayat pemakaian voucher: %w", err)
		}
		if used >= int64(v.PerUserLimit) {
			return nil, voucherModel.NewValidationError("kamu sudah memakai voucher ini")
		}
	}

	// Voucher milik seller tertentu hanya berlaku bila keranjang memang memuat
	// barang dari seller itu. Tanpa cek ini, voucher merchant bisa dipakai
	// untuk produk orang lain.
	if v.MerchantID != nil && *v.MerchantID != "" {
		if err := s.assertCartContainsMerchant(items, *v.MerchantID); err != nil {
			return nil, err
		}
	}

	return &voucherRedemption{Voucher: v, Discount: v.ComputeDiscount(amount)}, nil
}

func (s *service) assertCartContainsMerchant(items []order.CheckoutItem, merchantID string) error {
	variantIDs := make([]string, 0, len(items))
	for _, item := range items {
		variantIDs = append(variantIDs, item.VariantID)
	}

	merchantIDs, err := s.productVariantRepo.FindMerchantIDsByVariantIDs(variantIDs)
	if err != nil {
		return fmt.Errorf("gagal memverifikasi pemilik produk: %w", err)
	}
	for _, id := range merchantIDs {
		if id == merchantID {
			return nil
		}
	}
	return voucherModel.NewValidationError("voucher ini hanya berlaku untuk produk dari seller terkait")
}

func calculateDiscount(redemption *voucherRedemption, amount float64) float64 {
	if redemption == nil {
		return 0
	}
	return redemption.Discount
}

// redeemVoucher mengklaim voucher di dalam transaksi checkout. Kuota dinaikkan
// dengan syarat di dalam statement UPDATE, sehingga dua checkout bersamaan tidak
// bisa sama-sama melewati batas kuota.
func (s *service) redeemVoucher(tx *gorm.DB, redemption *voucherRedemption, orderID, userID string, discount float64) error {
	txVoucherRepo := s.voucherRepo.WithTx(tx)

	affected, err := txVoucherRepo.ConsumeQuota(redemption.Voucher.ID)
	if err != nil {
		return fmt.Errorf("gagal klaim kuota voucher: %w", err)
	}
	if affected == 0 {
		return voucherModel.NewValidationError("kuota voucher sudah habis")
	}

	return txVoucherRepo.CreateRedemption(voucherModel.VoucherRedemption{
		ID:             uuid.NewString(),
		VoucherID:      redemption.Voucher.ID,
		OrderID:        orderID,
		UserID:         userID,
		DiscountAmount: discount,
		CreatedAt:      time.Now(),
	})
}
