package content_page

import (
	modelContentPage "website-api/model/content-page"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Take(selectParams []string, conditions *modelContentPage.CmsPage) (cmsPage modelContentPage.CmsPage, err error)
		FaqList(reqQuerry *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error)
		FaqListAdmin(reqQuerry *modelContentPage.FaqListReqQuery) (resData []modelContentPage.CmsFaq, count int64, err error)
		FaqFindByID(id string) (modelContentPage.CmsFaq, error)
		FaqCreate(faq *modelContentPage.CmsFaq) error
		FaqUpdate(id string, values map[string]any) error
		FaqDelete(id string) error
		SeedCmsPage() error
		SeedCmsFaq() error
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
