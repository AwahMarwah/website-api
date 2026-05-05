package content_page

import modelContentPage "website-api/model/content-page"

func (r *repo) FaqList(reqQuerry *modelContentPage.FaqListReqQuery) (resData []modelContentPage.FaqListResponse, count int64, err error) {
	resData = make([]modelContentPage.FaqListResponse, 0)
	return resData, count, r.db.Debug().Model(&modelContentPage.CmsFaq{}).Count(&count).Limit(reqQuerry.Limit).Offset(reqQuerry.Offset).Scan(&resData).Error
}
