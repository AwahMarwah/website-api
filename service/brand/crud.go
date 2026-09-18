package brand

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	brandModel "website-api/model/brand"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *service) Create(req *brandModel.CreateBrandReq) (int, error) {
	b := brandModel.Brand{
		Id:      uuid.NewString(),
		Name:    req.Name,
		Slug:    req.Slug,
		LogoUrl: req.LogoUrl,
	}
	if err := s.brandRepo.Create(b); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membuat brand: %w", err)
	}
	return http.StatusCreated, nil
}

func (s *service) Update(id string, req *brandModel.UpdateBrandReq) (int, error) {
	if _, err := s.brandRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("brand not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil brand: %w", err)
	}

	values := map[string]any{}
	if req.Name != "" {
		values["name"] = req.Name
	}
	if req.Slug != "" {
		values["slug"] = req.Slug
	}
	values["logo_url"] = req.LogoUrl
	values["updated_at"] = time.Now()

	if err := s.brandRepo.Update(id, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui brand: %w", err)
	}
	return http.StatusOK, nil
}

func (s *service) Delete(id string) (int, error) {
	if s.brandRepo.Delete(id) != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus brand")
	}
	return http.StatusOK, nil
}
