package category

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"website-api/model/category"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *service) Create(req *category.CreateCategoryReq) (int, error) {
	cat := category.Category{
		Id:        uuid.NewString(),
		Name:      req.Name,
		Slug:      req.Slug,
		ParentId:  req.ParentId,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if err := s.categoryRepo.Create(cat); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membuat kategori: %w", err)
	}
	return http.StatusCreated, nil
}

func (s *service) Update(id string, req *category.UpdateCategoryReq) (int, error) {
	if _, err := s.categoryRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("category not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil kategori: %w", err)
	}

	values := map[string]any{
		"updated_at": time.Now().Format(time.RFC3339),
	}
	if req.Name != "" {
		values["name"] = req.Name
	}
	if req.Slug != "" {
		values["slug"] = req.Slug
	}
	if req.ParentId != "" {
		values["parent_id"] = req.ParentId
	}
	if err := s.categoryRepo.Update(id, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui kategori: %w", err)
	}
	return http.StatusOK, nil
}

func (s *service) Delete(id string) (int, error) {
	if err := s.categoryRepo.Delete(id); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus kategori: %w", err)
	}
	return http.StatusOK, nil
}