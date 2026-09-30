package merchant

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"website-api/common"
	productModel "website-api/model/product"
	product_variant "website-api/model/product-variant"
)

// resolveMerchantID menolerolkan userID menjadi merchant yang aktif.
//
// Dipakai seluruh endpoint seller. merchant_id SELALU diturunkan dari session dan
// tidak pernah dari body request: kalau request yang dipercaya, seller bisa menulis
// produk atas nama seller lain.
func (s *service) resolveMerchantID(userID string) (merchantID string, statusCode int, err error) {
	if userID == "" {
		return "", http.StatusUnauthorized, fmt.Errorf("user ID tidak ditemukan dalam context")
	}

	m, err := s.merchantRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || m.ID == "" {
			return "", http.StatusForbidden, fmt.Errorf("%s", common.MerchantNotRegistered)
		}
		return "", http.StatusInternalServerError, fmt.Errorf("gagal mengambil data merchant: %w", err)
	}
	if !m.IsActive {
		return "", http.StatusForbidden, fmt.Errorf("%s", common.MerchantNotApproved)
	}
	return m.ID, http.StatusOK, nil
}

// assertProductOwned memastikan produk benar-benar milik seller yang sedang login.
func (s *service) assertProductOwned(productID, merchantID string) (productModel.Product, int, error) {
	p, err := s.productRepo.FindOwnedByID(productID, merchantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p, http.StatusNotFound, fmt.Errorf("%s", common.MerchantForbidden)
		}
		return p, http.StatusInternalServerError, fmt.Errorf("gagal mengambil produk: %w", err)
	}
	return p, http.StatusOK, nil
}

// assertVariantOwned memastikan variant milik produk milik seller tersebut.
// Dilakukan lewat dua langkah karena variant tidak menyimpan merchant_id.
func (s *service) assertVariantOwned(variantID, merchantID string) (product_variant.ProductVariant, int, error) {
	variant, err := s.productVariantRepo.FindByID(variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return variant, http.StatusNotFound, fmt.Errorf("variant not found")
		}
		return variant, http.StatusInternalServerError, fmt.Errorf("gagal mengambil variant: %w", err)
	}
	if _, statusCode, err := s.assertProductOwned(variant.ProductID, merchantID); err != nil {
		return variant, statusCode, err
	}
	return variant, http.StatusOK, nil
}
