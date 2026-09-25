package content_page

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"website-api/cache"
	modelContentPage "website-api/model/content-page"
	content_page "website-api/repository/content-page"

	"gorm.io/gorm"
)

type (
	IService interface {
		Detail(reqPath *modelContentPage.ReqPath) (resData modelContentPage.DetailResponse, statusCode int, err error)
		FaqList(reqQuery *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error)
		FaqListAdmin(reqQuery *modelContentPage.FaqListReqQuery) (resData []modelContentPage.CmsFaq, count int64, err error)
		FaqCreate(req *modelContentPage.FaqCreateRequest) (int, error)
		FaqUpdate(id string, req *modelContentPage.FaqUpdateRequest) (int, error)
		FaqDelete(id string) (int, error)
		Seed() (err error)
		SeedCmsFaq() error
	}

	service struct {
		contentPageRepo content_page.IRepo
		redis           cache.Cache
	}
)

func NewService(contentPageRepo content_page.IRepo, redis cache.Cache) IService {
	return &service{
		contentPageRepo: contentPageRepo,
		redis:           redis,
	}
}

func (s *service) FaqListAdmin(reqQuery *modelContentPage.FaqListReqQuery) ([]modelContentPage.CmsFaq, int64, error) {
	return s.contentPageRepo.FaqListAdmin(reqQuery)
}

func (s *service) FaqCreate(req *modelContentPage.FaqCreateRequest) (int, error) {
	faq := &modelContentPage.CmsFaq{
		Question:  req.Question,
		Answer:    req.Answer,
		OrderNo:   req.OrderNo,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.contentPageRepo.FaqCreate(faq); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal membuat FAQ: %w", err)
	}
	s.redis.DeleteByPattern("content-page:*")
	return http.StatusCreated, nil
}

func (s *service) FaqUpdate(id string, req *modelContentPage.FaqUpdateRequest) (int, error) {
	existing, err := s.contentPageRepo.FaqFindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("FAQ not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil FAQ: %w", err)
	}

	values := map[string]any{"updated_at": time.Now()}
	if req.Question != "" {
		values["question"] = req.Question
	}
	if req.Answer != "" {
		values["answer"] = req.Answer
	}
	if req.OrderNo != nil {
		values["order_no"] = *req.OrderNo
	}
	if req.IsActive != nil {
		values["is_active"] = *req.IsActive
	}

	if err := s.contentPageRepo.FaqUpdate(existing.Id, values); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal memperbarui FAQ: %w", err)
	}
	s.redis.DeleteByPattern("content-page:*")
	return http.StatusOK, nil
}

func (s *service) FaqDelete(id string) (int, error) {
	if _, err := s.contentPageRepo.FaqFindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, fmt.Errorf("FAQ not found")
		}
		return http.StatusInternalServerError, fmt.Errorf("gagal mengambil FAQ: %w", err)
	}
	if err := s.contentPageRepo.FaqDelete(id); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("gagal menghapus FAQ: %w", err)
	}
	s.redis.DeleteByPattern("content-page:*")
	return http.StatusOK, nil
}

func (s *service) SeedCmsFaq() error {
	return s.contentPageRepo.SeedCmsFaq()
}