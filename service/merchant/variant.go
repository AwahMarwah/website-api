package merchant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	merchantModel "website-api/model/merchant"
	product_variant "website-api/model/product-variant"
)

// ListVariants mengembalikan seluruh variant milik satu produk milik seller.
func (s *service) ListVariants(userID, productID string) ([]merchantModel.SellerVariantResponse, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return nil, statusCode, err
	}
	if _, statusCode, err := s.assertProductOwned(productID, merchantID); err != nil {
		return nil, statusCode, err
	}

	variants, err := s.productVariantRepo.FindByProductID(productID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("gagal mengambil variant: %w", err)
	}

	res := make([]merchantModel.SellerVariantResponse, 0, len(variants))
	for _, v := range variants {
		res = append(res, mapVariant(v))
	}
	return res, http.StatusOK, nil
}

// CreateVariant menambah variant baru ke produk milik seller.
func (s *service) CreateVariant(userID, productID string, req *merchantModel.CreateSellerVariantReq) (string, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return "", statusCode, err
	}
	if _, statusCode, err := s.assertProductOwned(productID, merchantID); err != nil {
		return "", statusCode, err
	}

	variant := product_variant.ProductVariant{
		ID:          uuid.NewString(),
		ProductID:   productID,
		Sku:         req.Sku,
		VariantName: req.VariantName,
		Price:       req.Price,
		Stock:       req.Stock,
		Weight:      req.Weight,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}

	if err := s.productVariantRepo.Create(variant); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal membuat variant: %w", err)
	}
	return variant.ID, http.StatusCreated, nil
}

// UpdateVariant memperbarui variant milik seller.
//
// Stok TIDAK boleh diubah ke nilai yang lebih kecil dari stok yang sedang teralokasi
// pada order yang belum selesai; seller memanipulasi angka ini bisa membuat pesanan
// yang sudah dibayar tidak dapat dipenuhi.
func (s *service) UpdateVariant(userID, variantID string, req *merchantModel.UpdateSellerVariantReq) (int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return statusCode, err
	}
	variant, statusCode, err := s.assertVariantOwned(variantID, merchantID)
	if err != nil {
		return statusCode, err
	}

	values := map[string]any{}
	if req.Sku != nil {
		values["sku"] = *req.Sku
	}
	if req.VariantName != nil {
		values["variant_name"] = *req.VariantName
	}
	if req.Price != nil {
		values["price"] = *req.Price
	}
	if req.Weight != nil {
		values["weight"] = *req.Weight
	}
	if req.IsActive != nil {
		values["is_active"] = *req.IsActive
	}
	if req.Stock != nil {
		// Stok hanya boleh bertambah dari sisi seller. Pengurangan stok selling
		// dilakukan sistem saat pesanan dibuat, bukan input manual seller.
		if *req.Stock < variant.Stock {
			return http.StatusBadRequest, fmt.Errorf("stok tidak bisa dikurangi manual, hanya bisa ditambah")
		}
		values["stock"] = *req.Stock
	}

	if len(values) == 0 {
		return http.StatusOK, nil
	}
	if err := s.productVariantRepo.UpdateFields(variantID, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui variant: %w", err)
	}
	return http.StatusOK, nil
}

// DeleteVariant menonaktifkan variant. Barang yang sudah pernah dipesan tidak dihapus
// supaya histori order tetap utuh.
func (s *service) DeleteVariant(userID, variantID string) (int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return statusCode, err
	}
	if _, statusCode, err := s.assertVariantOwned(variantID, merchantID); err != nil {
		return statusCode, err
	}

	if err := s.productVariantRepo.UpdateFields(variantID, map[string]any{"is_active": false}); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menonaktifkan variant: %w", err)
	}
	return http.StatusOK, nil
}
