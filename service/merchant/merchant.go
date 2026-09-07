package merchant

import (
	"errors"
	"fmt"
	"net/http"
	merchantModel "website-api/model/merchant"
	"website-api/repository/merchant"

	"gorm.io/gorm"
)

type (
	IService interface {
		List() ([]merchantModel.MerchantResponse, int, error)
		Detail(id string) (merchantModel.MerchantResponse, int, error)
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

func mapToResponse(m merchantModel.Merchant) merchantModel.MerchantResponse {
	return merchantModel.MerchantResponse{
		ID:            m.ID,
		Name:          m.Name,
		Slug:          m.Slug,
		DestinationID: m.DestinationID,
		CityID:        m.CityID,
		Address:       m.Address,
		IsActive:      m.IsActive,
	}
}