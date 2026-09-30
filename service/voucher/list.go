package voucher

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"website-api/library/pagination"
	voucherModel "website-api/model/voucher"
)

func (s *service) List(reqQuery *voucherModel.ListVoucherReqQuery) ([]voucherModel.VoucherResponse, int64, int, error) {
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	vouchers, count, err := s.voucherRepo.List(reqQuery)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar voucher: %w", err)
	}

	res := make([]voucherModel.VoucherResponse, 0, len(vouchers))
	for _, v := range vouchers {
		res = append(res, mapToResponse(v))
	}
	return res, count, http.StatusOK, nil
}

func (s *service) Detail(id string) (voucherModel.VoucherResponse, int, error) {
	v, err := s.voucherRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return voucherModel.VoucherResponse{}, http.StatusNotFound, fmt.Errorf("voucher not found")
		}
		return voucherModel.VoucherResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengambil voucher: %w", err)
	}
	return mapToResponse(v), http.StatusOK, nil
}

func mapToResponse(v voucherModel.Voucher) voucherModel.VoucherResponse {
	res := voucherModel.VoucherResponse{
		ID:            v.ID,
		Code:          v.Code,
		Description:   v.Description,
		Type:          v.Type,
		Value:         v.Value,
		MaxDiscount:   v.MaxDiscount,
		MinSpend:      v.MinSpend,
		Quota:         v.Quota,
		UsedCount:     v.UsedCount,
		PerUserLimit:  v.PerUserLimit,
		StartsAt:      v.StartsAt,
		EndsAt:        v.EndsAt,
		MerchantID:    v.MerchantID,
		IsActive:      v.IsActive,
		CreatedAt:     v.CreatedAt,
		RemainingQuota: nil,
	}
	if v.Quota != nil {
		remaining := *v.Quota - v.UsedCount
		if remaining < 0 {
			remaining = 0
		}
		res.RemainingQuota = &remaining
	}
	return res
}
