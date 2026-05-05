package content_page

import (
	modelContentPage "website-api/model/content-page"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		Take(selectParams []string, conditions *modelContentPage.CmsPage) (cmsPage modelContentPage.CmsPage, err error)
		FaqList(reqQuerry *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error)
		SeedCmsPage() error
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) IRepo {
	return &repo{db: db}
}
