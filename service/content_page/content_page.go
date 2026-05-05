package content_page

import (
	"website-api/cache"
	modelContentPage "website-api/model/content-page"
	content_page "website-api/repository/content-page"
)

type (
	IService interface {
		Detail(reqPath *modelContentPage.ReqPath) (resData modelContentPage.DetailResponse, statusCode int, err error)
		FaqList(reqQuery *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error)
		Seed() (err error)
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
