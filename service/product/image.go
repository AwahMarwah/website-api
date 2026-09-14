package product

import (
	"fmt"
	"net/http"
	productModel "website-api/model/product"
)

func (s *service) SetPrimaryImage(productID, imageID string) (int, error) {
	if err := s.productRepo.SetPrimary(productID, imageID); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengatur gambar utama: %w", err)
	}
	return http.StatusOK, nil
}

func (s *service) DeleteImage(productID, imageID string) (int, error) {
	// Cek apakah gambar ini primary
	images, err := s.productRepo.FindImagesByProductID(productID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil gambar: %w", err)
	}

	var deletedPrimary bool
	for _, img := range images {
		if img.ID == imageID && img.IsPrimary {
			deletedPrimary = true
			break
		}
	}

	if err := s.productRepo.DeleteImage(imageID); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus gambar: %w", err)
	}

	// Jika yang dihapus adalah primary, promosikan gambar pertama tersisa
	if deletedPrimary {
		remaining, _ := s.productRepo.FindImagesByProductID(productID)
		if len(remaining) > 0 {
			s.productRepo.SetImagePrimaryFalse(imageID) // just in case
			_ = s.productRepo.SetPrimary(productID, remaining[0].ID) // reset: all false + first true
		}
	}

	return http.StatusOK, nil
}

func (s *service) ReorderImages(productID string, orders []productModel.ImageSortOrder) (int, error) {
	if err := s.productRepo.UpdateSortOrders(productID, orders); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal mengubah urutan gambar: %w", err)
	}
	return http.StatusOK, nil
}
