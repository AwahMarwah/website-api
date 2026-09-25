package category

import (
	"errors"
	"fmt"
	"net/http"
	"website-api/model/category"

	"gorm.io/gorm"
)

func (s *service) GetCategory(reqQuery *category.FilterCategory) (resData []*category.ListCategoryResponse, count int64, err error) {
	resData, count, err = s.categoryRepo.GetCategory(reqQuery)
	if err != nil {
		return resData, count, err
	}
	return
}

func (s *service) GetCategoryBySlug(slug string) (category.Category, int, error) {
	cat, err := s.categoryRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return category.Category{}, http.StatusNotFound, fmt.Errorf("kategori tidak ditemukan")
		}
		return category.Category{}, http.StatusInternalServerError, fmt.Errorf("gagal mengambil kategori: %w", err)
	}
	return cat, http.StatusOK, nil
}
