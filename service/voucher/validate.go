package voucher

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	voucherModel "website-api/model/voucher"
)

// validateVoucher memeriksa seluruh syarat voucher terhadap konteks pemesanan.
// Syarat non-database delegate ke model, syarat database dicek di sini.
func (s *service) validateVoucher(code, userID string, amount float64) (float64, error) {
	v, err := s.voucherRepo.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, voucherModel.NewValidationError("kode voucher tidak ditemukan")
		}
		return 0, fmt.Errorf("gagal mengambil voucher: %w", err)
	}

	if err := v.Validate(time.Now(), amount); err != nil {
		return 0, err
	}

	if v.PerUserLimit > 0 {
		used, err := s.voucherRepo.CountUserRedemptions(v.ID, userID)
		if err != nil {
			return 0, fmt.Errorf("gagal memeriksa riwayat pemakaian voucher: %w", err)
		}
		if used >= int64(v.PerUserLimit) {
			return 0, voucherModel.NewValidationError("kamu sudah memakai voucher ini")
		}
	}

	return v.ComputeDiscount(amount), nil
}
