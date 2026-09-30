package merchant

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	merchantModel "website-api/model/merchant"
)

// GetProduct mengembalikan detail satu produk milik seller.
func (s *service) GetProduct(userID, productID string) (merchantModel.SellerProductResponse, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return merchantModel.SellerProductResponse{}, statusCode, err
	}

	p, statusCode, err := s.assertProductOwned(productID, merchantID)
	if err != nil {
		return merchantModel.SellerProductResponse{}, statusCode, err
	}

	variants, err := s.productVariantRepo.FindByProductID(p.Id)
	if err != nil {
		return merchantModel.SellerProductResponse{}, http.StatusInternalServerError,
			fmt.Errorf("gagal mengambil variant: %w", err)
	}

	res := merchantModel.SellerProductResponse{
		ID:          p.Id,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		BrandId:     p.BrandId,
		Status:      p.Status,
		BasePrice:   p.BasePrice,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Variants:    make([]merchantModel.SellerVariantResponse, 0, len(variants)),
	}
	for _, v := range variants {
		res.Variants = append(res.Variants, mapVariant(v))
	}

	if images, err := s.productRepo.FindImagesByProductID(p.Id); err == nil {
		for _, img := range images {
			if img.IsPrimary {
				res.Thumbnail = img.ImageURL
				break
			}
		}
	}
	return res, http.StatusOK, nil
}

// UpdateProduct memperbarui produk milik seller. Field yang tidak dikirim tidak diubah.
func (s *service) UpdateProduct(userID, productID string, req *merchantModel.UpdateSellerProductReq) (int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return statusCode, err
	}
	if _, statusCode, err := s.assertProductOwned(productID, merchantID); err != nil {
		return statusCode, err
	}

	values := map[string]any{"updated_at": time.Now()}
	if req.Name != "" {
		values["name"] = req.Name
	}
	if req.Slug != "" {
		values["slug"] = req.Slug
	}
	if req.Description != "" {
		values["description"] = req.Description
	}
	if req.BrandId != nil {
		values["brand_id"] = *req.BrandId
	}
	if req.BasePrice != nil {
		values["base_price"] = *req.BasePrice
	}
	if req.WeightGram != nil {
		values["weight_gram"] = *req.WeightGram
	}

	err = s.txManager.Execute(func(tx *gorm.DB) error {
		txProductRepo := s.productRepo.WithTx(tx)
		if err := txProductRepo.UpdateProduct(productID, values); err != nil {
			return err
		}
		if req.Categories != nil {
			return txProductRepo.RebindCategories(productID, req.Categories)
		}
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui produk: %w", err)
	}
	return http.StatusOK, nil
}

// DeleteProduct menonaktifkan produk milik seller.
//
// Soft delete: produk yang pernah dipesan tidak boleh hilang dari histori order.
// Seller juga tidak boleh menonaktifkan produk yang masih ada di order yang belum selesai.
func (s *service) DeleteProduct(userID, productID string) (int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return statusCode, err
	}
	if _, statusCode, err := s.assertProductOwned(productID, merchantID); err != nil {
		return statusCode, err
	}

	if err := s.productRepo.SetInactive(productID); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menonaktifkan produk: %w", err)
	}
	return http.StatusOK, nil
}
