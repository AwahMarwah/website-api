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
		Id       string
		Question string
		Answer   string
		OrderNo  int
		IsActive bool
	}
)

func (cmsPage *CmsPage) BeforeCreate(*gorm.DB) error {
	cmsPage.ID = uuid.New().String()
	return nil
}
