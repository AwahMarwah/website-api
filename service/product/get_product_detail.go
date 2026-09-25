package product

import (
	"errors"
	"fmt"
	productModel "website-api/model/product"

	"gorm.io/gorm"
)

func (s *service) GetProductDetail(id string) (resData productModel.ProductDetailResponse, err error) {
	resData, err = s.productRepo.FindDetailByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, fmt.Errorf("product not found")
		}
		return resData, fmt.Errorf("gagal mengambil detail produk: %w", err)
	}

	variants, err := s.productVariantRepo.FindByProductID(id)
	if err != nil {
		return resData, fmt.Errorf("gagal mengambil varian produk: %w", err)
	}

	for _, v := range variants {
		resData.Variants = append(resData.Variants, productModel.VariantResponse{
			ID:          v.ID,
			Sku:         v.Sku,
			VariantName: v.VariantName,
			Price:       v.Price,
			Stock:       v.Stock,
			Weight:      v.Weight,
		})
	}

	// Load gallery images (eager loading via Preload)
	product, err := s.productRepo.FindDetailCollections(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resData, fmt.Errorf("product not found")
		}
		return resData, fmt.Errorf("gagal mengambil gambar produk: %w", err)
	}
	for _, img := range product.Images {
		resData.Images = append(resData.Images, productModel.ImageResponse{
			ID:        img.ID,
			URL:       img.ImageURL,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
		})
	}

	return resData, nil
}
