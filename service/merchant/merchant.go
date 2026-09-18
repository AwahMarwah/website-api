package merchant

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	merchantModel "website-api/model/merchant"
	"website-api/repository/merchant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IService interface {
		List() ([]merchantModel.MerchantResponse, int, error)
		Detail(id string) (merchantModel.MerchantResponse, int, error)
		Register(req *merchantModel.CreateMerchantReq, userID string) (int, error)
		GetMy(userID string) (merchantModel.MerchantResponse, int, error)
		UpdateMy(userID string, req *merchantModel.UpdateMerchantReq) (int, error)
		Approve(id string) (int, error)
	}

	service struct {
		merchantRepo merchant.IRepo
	}
)

func NewService(merchantRepo merchant.IRepo) IService {
	return &service{merchantRepo: merchantRepo}
}

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

func (s *service) Register(req *merchantModel.CreateMerchantReq, userID string) (int, error) {
	// Cek apakah user sudah punya merchant
	if existing, err := s.merchantRepo.FindByUserID(userID); err == nil && existing.ID != "" {
		return http.StatusConflict, fmt.Errorf("akun ini sudah terdaftar sebagai merchant")
	}

	m := merchantModel.Merchant{
		ID:            uuid.NewString(),
		Name:          req.Name,
		Slug:          req.Slug,
		DestinationID: req.DestinationID,
		CityID:        req.CityID,
		Address:       req.Address,
		IsActive:      false, // menunggu approval admin
		UserID:        userID,
		CreatedAt:     time.Now(),
	}

	if err := s.merchantRepo.Create(m); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mendaftarkan merchant: %w", err)
	}
	return http.StatusCreated, nil
}

func (s *service) GetMy(userID string) (merchantModel.MerchantResponse, int, error) {
	m, err := s.merchantRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merchantModel.MerchantResponse{}, http.StatusNotFound, fmt.Errorf("merchant tidak ditemukan untuk akun ini")
		}
		return merchantModel.MerchantResponse{}, http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}
	return mapToResponse(m), http.StatusOK, nil
}

func (s *service) UpdateMy(userID string, req *merchantModel.UpdateMerchantReq) (int, error) {
	existing, err := s.merchantRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("merchant tidak ditemukan")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil merchant: %w", err)
	}

	values := map[string]any{}
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
	values["updated_at"] = time.Now()

	if err := s.merchantRepo.Update(existing.ID, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui merchant: %w", err)
	}
	return http.StatusOK, nil
}

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
		ID:            m.ID,
		Name:          m.Name,
		Slug:          m.Slug,
		DestinationID: m.DestinationID,
		CityID:        m.CityID,
		Address:       m.Address,
		IsActive:      m.IsActive,
		UserID:        m.UserID,
	}
}