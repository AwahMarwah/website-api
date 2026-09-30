package voucher

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"website-api/common"
	voucherModel "website-api/model/voucher"
)

func (s *service) Create(req *voucherModel.CreateVoucherReq) (string, int, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))

	if req.Type == common.DiscountTypePercent && req.Value > 100 {
		return "", 400, fmt.Errorf("nilai persentase voucher tidak boleh melebihi 100")
	}

	startsAt, endsAt, err := parseWindow(req.StartsAt, req.EndsAt)
	if err != nil {
		return "", 400, err
	}

	v := voucherModel.Voucher{
		ID:           uuid.NewString(),
		Code:         code,
		Description:  optionalString(req.Description),
		Type:         req.Type,
		Value:        req.Value,
		MaxDiscount:  req.MaxDiscount,
		MinSpend:     req.MinSpend,
		Quota:        req.Quota,
		UsedCount:    0,
		PerUserLimit: req.PerUserLimit,
		StartsAt:     startsAt,
		EndsAt:       endsAt,
		MerchantID:   req.MerchantID,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	if v.PerUserLimit <= 0 {
		v.PerUserLimit = 1
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		txRepo := s.voucherRepo.WithTx(tx)
		if err := txRepo.Create(v); err != nil {
			return err
		}
		if len(req.ProductIDs) > 0 {
			return txRepo.RebindProducts(v.ID, req.ProductIDs)
		}
		return nil
	})
	if err != nil {
		return "", 500, fmt.Errorf("gagal membuat voucher: %w", err)
	}

	return v.ID, 201, nil
}

func parseWindow(startsAt, endsAt *string) (*time.Time, *time.Time, error) {
	var start, end *time.Time

	if startsAt != nil && *startsAt != "" {
		parsed, err := time.Parse(time.RFC3339, *startsAt)
		if err != nil {
			return nil, nil, fmt.Errorf("starts_at tidak valid, gunakan format RFC3339")
		}
		start = &parsed
	}
	if endsAt != nil && *endsAt != "" {
		parsed, err := time.Parse(time.RFC3339, *endsAt)
		if err != nil {
			return nil, nil, fmt.Errorf("ends_at tidak valid, gunakan format RFC3339")
		}
		end = &parsed
	}
	if start != nil && end != nil && !end.After(*start) {
		return nil, nil, fmt.Errorf("ends_at harus setelah starts_at")
	}
	return start, end, nil
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
