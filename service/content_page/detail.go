package content_page

import (
	"net/http"
	modelContentPage "website-api/model/content-page"
)

func (s *service) Detail(reqPath *modelContentPage.ReqPath) (resData modelContentPage.DetailResponse, statusCode int, err error) {
	contentPage, err := s.contentPageRepo.Take([]string{"id", "slug", "title", "content", "status"}, &modelContentPage.CmsPage{Slug: reqPath.Slug})
	if err != nil {
		return resData, http.StatusInternalServerError, err
	}
	resData = modelContentPage.DetailResponse{
		ID:      contentPage.ID,
		Slug:    contentPage.Slug,
		Title:   contentPage.Title,
		Content: contentPage.Content,
		Status:  contentPage.Status,
	}
	return resData, statusCode, err
}
