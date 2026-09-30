package voucher

import (
	"fmt"

	voucherModel "website-api/model/voucher"
)

// Validate mengecek voucher tanpa mengubah apa pun, untuk pratinjau di frontend.
// Selalu mengembalikan HTTP 200 dengan flag valid: voucher yang tidak berlaku
// adalah jawaban yang wajar, bukan kegagalan endpoint.
func (s *service) Validate(code, userID string, amount float64) (voucherModel.ValidateVoucherResponse, error) {
	res := voucherModel.ValidateVoucherResponse{Code: code, Valid: false, FinalAmount: amount}

	discount, err := s.validateVoucher(code, userID, amount)
	if err != nil {
		if verr, ok := voucherModel.AsValidationError(err); ok {
			res.Message = verr.Message
			return res, nil
		}
		return res, fmt.Errorf("gagal memvalidasi voucher: %w", err)
	}

	res.Valid = true
	res.DiscountAmount = discount
	res.FinalAmount = amount - discount
	res.Message = "voucher berhasil digunakan"
	return res, nil
}
