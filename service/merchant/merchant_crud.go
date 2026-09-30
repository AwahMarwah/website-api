package merchant

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	merchantModel "website-api/model/merchant"
)

func (s *service) List() ([]merchantModel.MerchantResponse, int, error) {
	merchants, err := s.merchantRepo.FindAll()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}
	res := make([]merchantModel.MerchantResponse, 0, len(merchants))
	for _, m := range merchants {
		res = append(res, mapToResponse(m))
	}
	return res, http.StatusOK, nil
}

func (s *service) Detail(id string) (merchantModel.MerchantResponse, int, error) {
	m, err := s.merchantRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merchantModel.MerchantResponse{}, http.StatusNotFound, fmt.Errorf("merchant not found")
		}
		return merchantModel.MerchantResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}
	return mapToResponse(m), http.StatusOK, nil
}

func (s *service) GetMy(userID string) (merchantModel.MerchantResponse, int, error) {
	m, err := s.merchantRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || m.ID == "" {
			return merchantModel.MerchantResponse{}, http.StatusNotFound, fmt.Errorf("merchant tidak ditemukan untuk akun ini")
		}
		return merchantModel.MerchantResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}
	return mapToResponse(m), http.StatusOK, nil
}

// UpdateMy memperbarui profil merchant sendiri.
func (s *service) UpdateMy(userID string, req *merchantModel.UpdateMerchantReq) (int, error) {
	existing, err := s.merchantRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || existing.ID == "" {
			return http.StatusNotFound, fmt.Errorf("merchant tidak ditemukan")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}

	values := map[string]any{"updated_at": time.Now()}
	if req.Name != "" {
		values["name"] = req.Name
	}
	if req.Slug != "" {
		values["slug"] = req.Slug
	}
	if req.DestinationID != nil {
		values["destination_id"] = *req.DestinationID
	}
	if req.CityID != "" {
		values["city_id"] = req.CityID
	}
	if req.Address != "" {
		values["address"] = req.Address
	}
	// IsActive sengaja tidak bisa diubah dari sini: hanya admin yang boleh
	// menonaktifkan merchant lewat endpoint Approve.

	if err := s.merchantRepo.Update(existing.ID, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui merchant: %w", err)
	}
	return http.StatusOK, nil
}

// Approve mengaktifkan merchant. Hanya admin yang boleh memanggilnya.
func (s *service) Approve(id string) (int, error) {
	existing, err := s.merchantRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("merchant not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}
	if existing.IsActive {
		return http.StatusConflict, fmt.Errorf("merchant sudah aktif")
	}

	if err := s.merchantRepo.Update(id, map[string]any{"is_active": true}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menyetujui merchant: %w", err)
	}
	return http.StatusOK, nil
}

func mapToResponse(m merchantModel.Merchant) merchantModel.MerchantResponse {
	return merchantModel.MerchantResponse{
		ID:               m.ID,
		Name:             m.Name,
		Slug:             m.Slug,
		DestinationID:    m.DestinationID,
		CityID:           m.CityID,
		Address:          m.Address,
		IsActive:         m.IsActive,
		UserID:           m.UserID,
		CommissionRateBP: m.CommissionRateBP,
	}
}
