package content_page

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CmsPage struct {
		ID        string
		Slug      string
		Title     string
		Content   string
		Status    bool
		CreatedAt time.Time
		CreatedBy string
		UpdatedAt time.Time
		UpdatedBy string
	}

	CmsFaq struct {
		Id        string    `json:"id" gorm:"column:id"`
		Question  string    `json:"question" gorm:"column:question"`
		Answer    string    `json:"answer" gorm:"column:answer"`
		OrderNo   int       `json:"order_no" gorm:"column:order_no"`
		IsActive  bool      `json:"is_active" gorm:"column:is_active"`
		CreatedAt time.Time `json:"created_at"`
		CreatedBy string    `json:"created_by"`
		UpdatedAt time.Time `json:"updated_at"`
		UpdatedBy string    `json:"updated_by"`
	}

	FaqCreateRequest struct {
		Question string `binding:"required" json:"question"`
		Answer   string `binding:"required" json:"answer"`
		OrderNo  int    `json:"order_no"`
	}

	FaqUpdateRequest struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
		OrderNo  *int   `json:"order_no"`
		IsActive *bool  `json:"is_active"`
	}

	FaqReqPath struct {
		Id string `uri:"id" binding:"required"`
	}
)

func (cmsPage *CmsPage) BeforeCreate(*gorm.DB) error {
	cmsPage.ID = uuid.New().String()
	return nil
}

func (cmsFaq *CmsFaq) BeforeCreate(*gorm.DB) error {
	if cmsFaq.Id == "" {
		cmsFaq.Id = uuid.New().String()
	}
	return nil
}