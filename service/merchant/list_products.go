package merchant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"website-api/common"
	merchantModel "website-api/model/merchant"
	productModel "website-api/model/product"
	product_variant "website-api/model/product-variant"
)

// ListProducts mengembalikan daftar produk milik seller yang sedang login.
func (s *service) ListProducts(userID string, reqQuery *merchantModel.ListSellerProductReqQuery) ([]merchantModel.SellerProductResponse, int64, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return nil, 0, statusCode, err
	}

	products, count, err := s.productRepo.FindByMerchantID(
		merchantID, reqQuery.Status, reqQuery.Search, "", reqQuery.Sort,
		reqQuery.Limit, reqQuery.Offset,
	)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil daftar produk: %w", err)
	}

	res := make([]merchantModel.SellerProductResponse, 0, len(products))
	for _, p := range products {
		item := merchantModel.SellerProductResponse{
			ID:        p.Id,
			Name:      p.Name,
			Slug:      p.Slug,
			BrandId:   p.BrandId,
			Status:    p.Status,
			Thumbnail: p.Thumbnail,
		}
		variants, err := s.productVariantRepo.FindByProductID(p.Id)
		if err != nil {
			return nil, 0, http.StatusInternalServerError, fmt.Errorf("gagal mengambil variant produk: %w", err)
		}
		item.Variants = make([]merchantModel.SellerVariantResponse, 0, len(variants))
		for _, v := range variants {
			item.Variants = append(item.Variants, mapVariant(v))
		}
		res = append(res, item)
	}
	return res, count, http.StatusOK, nil
}

func mapVariant(v product_variant.ProductVariant) merchantModel.SellerVariantResponse {
	return merchantModel.SellerVariantResponse{
		ID:          v.ID,
		Sku:         v.Sku,
		VariantName: v.VariantName,
		Price:       v.Price,
		Stock:       v.Stock,
		Weight:      v.Weight,
		IsActive:    v.IsActive,
	}
}

// CreateProduct membuat produk baru milik seller.
//
// Status selalu draft: seller tidak boleh menerbitkan produk langsung ke katalog
// tanpa ditinjau admin.
func (s *service) CreateProduct(userID string, req *merchantModel.CreateSellerProductReq) (string, int, error) {
	merchantID, statusCode, err := s.resolveMerchantID(userID)
	if err != nil {
		return "", statusCode, err
	}

	product := productModel.Product{
		Id:          uuid.NewString(),
		BrandId:     req.BrandId,
		MerchantId:  merchantID,
		Sku:         req.Sku,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		BasePrice:   float32(req.BasePrice),
		Status:      common.ProductStatusDraft,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.productRepo.Create(product); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("gagal membuat produk: %w", err)
	}

	if len(req.Categories) > 0 {
		if err := s.productRepo.RebindCategories(product.Id, req.Categories); err != nil {
			return "", http.StatusInternalServerError, fmt.Errorf("gagal mengaitkan kategori: %w", err)
		}
	}
	return product.Id, http.StatusCreated, nil
}
