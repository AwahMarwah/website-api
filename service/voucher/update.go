package voucher

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	voucherModel "website-api/model/voucher"
)

func (s *service) Update(id string, req *voucherModel.UpdateVoucherReq) (int, error) {
	existing, err := s.voucherRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("voucher not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil voucher: %w", err)
	}

	values := map[string]any{"updated_at": time.Now()}

	if req.Description != nil {
		values["description"] = *req.Description
	}
	if req.Value != nil {
		if existing.Type == "PERCENT" && *req.Value > 100 {
			return http.StatusBadRequest, fmt.Errorf("nilai persentase voucher tidak boleh melebihi 100")
		}
		values["value"] = *req.Value
	}
	if req.MaxDiscount != nil {
		values["max_discount"] = *req.MaxDiscount
	}
	if req.MinSpend != nil {
		values["min_spend"] = *req.MinSpend
	}
	if req.Quota != nil {
		values["quota"] = *req.Quota
	}
	if req.PerUserLimit != nil {
		values["per_user_limit"] = *req.PerUserLimit
	}
	if req.IsActive != nil {
		values["is_active"] = *req.IsActive
	}

	// Periode hanya bisa diubah lewat kedua nilai sekaligus supaya tidak bisa
	// menghasilkan voucher dengan tanggal terbalik.
	if req.StartsAt != nil || req.EndsAt != nil {
		startStr, endStr := "", ""
		if req.StartsAt != nil {
			startStr = *req.StartsAt
		} else if existing.StartsAt != nil {
			startStr = existing.StartsAt.Format(time.RFC3339)
		}
		if req.EndsAt != nil {
			endStr = *req.EndsAt
		} else if existing.EndsAt != nil {
			endStr = existing.EndsAt.Format(time.RFC3339)
		}
		start, end, err := parseWindow(&startStr, &endStr)
		if err != nil {
			return http.StatusBadRequest, err
		}
		values["starts_at"] = start
		values["ends_at"] = end
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		txRepo := s.voucherRepo.WithTx(tx)
		if err := txRepo.Update(id, values); err != nil {
			return err
		}
		if req.ProductIDs != nil {
			return txRepo.RebindProducts(id, req.ProductIDs)
		}
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui voucher: %w", err)
	}
	return http.StatusOK, nil
}

// Delete menonaktifkan voucher, tidak menghapusnya.
// Order lama yang sudah memakainya harus tetap bisa ditelusuri lewat redemptions.
func (s *service) Delete(id string) (int, error) {
	existing, err := s.voucherRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("voucher not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil voucher: %w", err)
	}
	if err := s.voucherRepo.Update(existing.ID, map[string]any{"is_active": false}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menonaktifkan voucher: %w", err)
	}
	return http.StatusOK, nil
}
