package product

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	productModel "website-api/model/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *service) CreateProduct(req *productModel.CreateProductReq) (int, error) {
	productID := uuid.NewString()
	err := s.txManager.Execute(func(tx *gorm.DB) error {
		txProductRepo := s.productRepo.WithTx(tx)

		p := productModel.Product{
			Id:          productID,
			Name:        req.Name,
			BrandId:     req.BrandId,
			MerchantId:  req.MerchantId,
			Sku:         req.Sku,
			Slug:        req.Slug,
			Description: req.Description,
			BasePrice:   float32(req.BasePrice),
			Status:      "active",
			CreatedAt:   time.Now(),
		}
		if err := txProductRepo.Create(p); err != nil {
			return err
		}

		if len(req.Categories) > 0 {
			if err := txProductRepo.RebindCategories(productID, req.Categories); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membuat produk: %w", err)
	}
	s.cache.DeleteByPattern("product:*")
	return http.StatusCreated, nil
}

func (s *service) UpdateProduct(id string, req *productModel.UpdateProductReq) (int, error) {
	if _, err := s.productRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("product not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil produk: %w", err)
	}

	err := s.txManager.Execute(func(tx *gorm.DB) error {
		txProductRepo := s.productRepo.WithTx(tx)

		values := map[string]any{"updated_at": time.Now()}
		if req.Name != "" {
			values["name"] = req.Name
		}
		if req.BrandId != "" {
			values["brand_id"] = req.BrandId
		}
		if req.MerchantId != "" {
			values["merchant_id"] = req.MerchantId
		}
		if req.Sku != "" {
			values["sku"] = req.Sku
		}
		if req.Slug != "" {
			values["slug"] = req.Slug
		}
		if req.Description != "" {
			values["description"] = req.Description
		}
		if req.BasePrice != 0 {
			values["base_price"] = req.BasePrice
		}
		if req.Status != "" {
			values["status"] = req.Status
		}

		if err := txProductRepo.UpdateProduct(id, values); err != nil {
			return err
		}

		if req.Categories != nil {
			if err := txProductRepo.RebindCategories(id, req.Categories); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui produk: %w", err)
	}
	s.cache.DeleteByPattern("product:*")
	return http.StatusOK, nil
}

func (s *service) DeleteProduct(id string) (int, error) {
	if err := s.productRepo.SetInactive(id); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menonaktifkan produk: %w", err)
	}
	s.cache.DeleteByPattern("product:*")
	return http.StatusOK, nil
}