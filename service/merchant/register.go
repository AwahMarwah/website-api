package merchant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	merchantModel "website-api/model/merchant"
)

// defaultCommissionRateBP adalah komisi default marketplace dalam basis points (10%).
// Nilai yang sama dipakai sebagai default kolom di migrasi.
const defaultCommissionRateBP = 1000

// Register mendaftarkan user sebagai merchant. Pendaftaran baru menghasilkan merchant
// dengan IsActive false; seller baru bisa mengelola produk setelah admin menyetujui.
func (s *service) Register(req *merchantModel.CreateMerchantReq, userID string) (int, error) {
	if existing, err := s.merchantRepo.FindByUserID(userID); err == nil && existing.ID != "" {
		return http.StatusConflict, fmt.Errorf("akun ini sudah terdaftar sebagai merchant")
	}

	m := merchantModel.Merchant{
		ID:               uuid.NewString(),
		Name:             req.Name,
		Slug:             req.Slug,
		DestinationID:    req.DestinationID,
		CityID:           req.CityID,
		Address:          req.Address,
		IsActive:         false, // menunggu approval admin
		UserID:           userID,
		CommissionRateBP: defaultCommissionRateBP,
		CreatedAt:        time.Now(),
	}

	if err := s.merchantRepo.Create(m); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mendaftarkan merchant: %w", err)
	}
	return http.StatusCreated, nil
}
